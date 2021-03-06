package mongo

import (
	"context"
	"time"

	mgo "github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"

	"github.com/tidepool-org/platform/blob"
	blobStoreStructured "github.com/tidepool-org/platform/blob/store/structured"
	"github.com/tidepool-org/platform/errors"
	"github.com/tidepool-org/platform/log"
	"github.com/tidepool-org/platform/page"
	"github.com/tidepool-org/platform/pointer"
	"github.com/tidepool-org/platform/request"
	storeStructuredMongo "github.com/tidepool-org/platform/store/structured/mongo"
	structureValidator "github.com/tidepool-org/platform/structure/validator"
	"github.com/tidepool-org/platform/user"
)

type Store struct {
	*storeStructuredMongo.Store
}

var (
	blobIndexes = map[string][]mgo.Index{
		"blobs": {
			{Key: []string{"id"}, Background: true, Unique: true},
			{Key: []string{"userId"}, Background: true},
			{Key: []string{"mediaType"}, Background: true},
			{Key: []string{"status"}, Background: true},
		},
	}
)

func NewStore(config *storeStructuredMongo.Config, logger log.Logger) (*Store, error) {
	if config != nil {
		config.Indexes = blobIndexes
	}
	store, err := storeStructuredMongo.NewStore(config, logger)
	if err != nil {
		return nil, err
	}

	return &Store{
		Store: store,
	}, nil
}

func (s *Store) NewSession() blobStoreStructured.Session {
	return s.newSession()
}

func (s *Store) newSession() *Session {
	return &Session{
		Session: s.Store.NewSession("blobs"),
	}
}

type Session struct {
	*storeStructuredMongo.Session
}

func (s *Session) List(ctx context.Context, userID string, filter *blob.Filter, pagination *page.Pagination) (blob.BlobArray, error) {
	ctx, logger := log.ContextAndLoggerWithFields(ctx, log.Fields{"userId": userID, "filter": filter, "pagination": pagination})

	if ctx == nil {
		return nil, errors.New("context is missing")
	}
	if userID == "" {
		return nil, errors.New("user id is missing")
	} else if !user.IsValidID(userID) {
		return nil, errors.New("user id is invalid")
	}
	if filter == nil {
		filter = blob.NewFilter()
	} else if err := structureValidator.New().Validate(filter); err != nil {
		return nil, errors.Wrap(err, "filter is invalid")
	}
	if pagination == nil {
		pagination = page.NewPagination()
	} else if err := structureValidator.New().Validate(pagination); err != nil {
		return nil, errors.Wrap(err, "pagination is invalid")
	}

	if s.IsClosed() {
		return nil, errors.New("session closed")
	}

	now := time.Now()

	var status []string
	if filter.Status != nil {
		status = *filter.Status
	} else {
		status = []string{blob.StatusAvailable}
	}

	result := blob.BlobArray{}
	query := bson.M{
		"userId": userID,
		"status": bson.M{
			"$in": status,
		},
		"deletedTime": bson.M{
			"$exists": false,
		},
	}
	if filter.MediaType != nil {
		query["mediaType"] = bson.M{
			"$in": *filter.MediaType,
		}
	}
	err := s.C().Find(query).Sort("-createdTime").Skip(pagination.Page * pagination.Size).Limit(pagination.Size).All(&result)
	if err != nil {
		logger.WithError(err).Error("Unable to list blobs")
		return nil, errors.Wrap(err, "unable to list blobs")
	}

	logger.WithFields(log.Fields{"count": len(result), "duration": time.Since(now) / time.Microsecond}).Debug("List")
	return result, nil
}

func (s *Session) Create(ctx context.Context, userID string, create *blobStoreStructured.Create) (*blob.Blob, error) {
	ctx, logger := log.ContextAndLoggerWithFields(ctx, log.Fields{"userId": userID, "create": create})

	if ctx == nil {
		return nil, errors.New("context is missing")
	}
	if userID == "" {
		return nil, errors.New("user id is missing")
	} else if !user.IsValidID(userID) {
		return nil, errors.New("user id is invalid")
	}
	if create == nil {
		return nil, errors.New("create is missing")
	} else if err := structureValidator.New().Validate(create); err != nil {
		return nil, errors.Wrap(err, "create is invalid")
	}

	if s.IsClosed() {
		return nil, errors.New("session closed")
	}

	now := time.Now()

	doc := &blob.Blob{
		UserID:      pointer.FromString(userID),
		MediaType:   create.MediaType,
		Status:      pointer.FromString(blob.StatusCreated),
		CreatedTime: pointer.FromTime(now.Truncate(time.Millisecond)),
		Revision:    pointer.FromInt(0),
	}

	var id string
	var err error
	for retry := 0; retry < 3; retry++ {
		id = blob.NewID()
		logger = logger.WithField("id", id)

		doc.ID = pointer.FromString(id)
		if err = s.C().Insert(doc); mgo.IsDup(err) {
			logger.WithError(err).Error("Duplicate blob id")
		} else {
			break
		}
	}
	if err != nil {
		logger.WithError(err).Error("Unable to create blob")
		return nil, errors.Wrap(err, "unable to create blob")
	}

	result, err := s.get(logger, id, nil)
	if err != nil {
		return nil, err
	}

	logger.WithField("duration", time.Since(now)/time.Microsecond).Debug("Create")
	return result, nil
}

func (s *Session) DeleteAll(ctx context.Context, userID string) (bool, error) {
	ctx, logger := log.ContextAndLoggerWithField(ctx, "userId", userID)

	if ctx == nil {
		return false, errors.New("context is missing")
	}
	if userID == "" {
		return false, errors.New("user id is missing")
	} else if !user.IsValidID(userID) {
		return false, errors.New("user id is invalid")
	}

	if s.IsClosed() {
		return false, errors.New("session closed")
	}

	now := time.Now()

	query := bson.M{
		"userId": userID,
	}
	set := bson.M{
		"modifiedTime": now.Truncate(time.Millisecond),
		"deletedTime":  now.Truncate(time.Millisecond),
	}
	unset := bson.M{}
	changeInfo, err := s.C().UpdateAll(query, s.ConstructUpdate(set, unset))
	if err != nil {
		logger.WithError(err).Error("Unable to delete all blobs")
		return false, errors.Wrap(err, "unable to delete all blobs")
	}

	logger.WithFields(log.Fields{"changeInfo": changeInfo, "duration": time.Since(now) / time.Microsecond}).Debug("DeleteAll")
	return changeInfo.Updated > 0, nil
}

func (s *Session) DestroyAll(ctx context.Context, userID string) (bool, error) {
	ctx, logger := log.ContextAndLoggerWithField(ctx, "userId", userID)

	if ctx == nil {
		return false, errors.New("context is missing")
	}
	if userID == "" {
		return false, errors.New("user id is missing")
	} else if !user.IsValidID(userID) {
		return false, errors.New("user id is invalid")
	}

	if s.IsClosed() {
		return false, errors.New("session closed")
	}

	now := time.Now()

	query := bson.M{
		"userId": userID,
	}
	changeInfo, err := s.C().RemoveAll(query)
	if err != nil {
		logger.WithError(err).Error("Unable to destroy all blobs")
		return false, errors.Wrap(err, "unable to destroy all blobs")
	}

	logger.WithFields(log.Fields{"changeInfo": changeInfo, "duration": time.Since(now) / time.Microsecond}).Debug("DestroyAll")
	return changeInfo.Removed > 0, nil
}

func (s *Session) Get(ctx context.Context, id string, condition *request.Condition) (*blob.Blob, error) {
	ctx, logger := log.ContextAndLoggerWithFields(ctx, log.Fields{"id": id, "condition": condition})

	if ctx == nil {
		return nil, errors.New("context is missing")
	}
	if id == "" {
		return nil, errors.New("id is missing")
	} else if !blob.IsValidID(id) {
		return nil, errors.New("id is invalid")
	}
	if condition == nil {
		condition = request.NewCondition()
	} else if err := structureValidator.New().Validate(condition); err != nil {
		return nil, errors.Wrap(err, "condition is invalid")
	}

	if s.IsClosed() {
		return nil, errors.New("session closed")
	}

	now := time.Now()

	result, err := s.get(logger, id, condition, storeStructuredMongo.NotDeleted)
	if err != nil {
		return nil, err
	}

	logger.WithField("duration", time.Since(now)/time.Microsecond).Debug("Get")
	return result, nil
}

func (s *Session) Update(ctx context.Context, id string, condition *request.Condition, update *blobStoreStructured.Update) (*blob.Blob, error) {
	ctx, logger := log.ContextAndLoggerWithFields(ctx, log.Fields{"id": id, "condition": condition, "update": update})

	if ctx == nil {
		return nil, errors.New("context is missing")
	}
	if id == "" {
		return nil, errors.New("id is missing")
	} else if !blob.IsValidID(id) {
		return nil, errors.New("id is invalid")
	}
	if condition == nil {
		condition = request.NewCondition()
	} else if err := structureValidator.New().Validate(condition); err != nil {
		return nil, errors.Wrap(err, "condition is invalid")
	}
	if update == nil {
		return nil, errors.New("update is missing")
	} else if err := structureValidator.New().Validate(update); err != nil {
		return nil, errors.Wrap(err, "update is invalid")
	}

	if s.IsClosed() {
		return nil, errors.New("session closed")
	}

	now := time.Now()

	if !update.IsEmpty() {
		query := bson.M{
			"id": id,
			"deletedTime": bson.M{
				"$exists": false,
			},
		}
		if condition.Revision != nil {
			query["revision"] = *condition.Revision
		}
		set := bson.M{
			"modifiedTime": now.Truncate(time.Millisecond),
		}
		unset := bson.M{}
		if update.MediaType != nil {
			set["mediaType"] = *update.MediaType
		}
		if update.DigestMD5 != nil {
			set["digestMD5"] = *update.DigestMD5
		}
		if update.Size != nil {
			set["size"] = *update.Size
		}
		if update.Status != nil {
			set["status"] = *update.Status
		}
		changeInfo, err := s.C().UpdateAll(query, s.ConstructUpdate(set, unset))
		if err != nil {
			logger.WithError(err).Error("Unable to update blob")
			return nil, errors.Wrap(err, "unable to update blob")
		} else if changeInfo.Matched > 0 {
			condition = nil
		} else {
			update = nil
		}

		logger = logger.WithField("changeInfo", changeInfo)
	}

	var result *blob.Blob
	if update != nil {
		var err error
		if result, err = s.get(logger, id, condition); err != nil {
			return nil, err
		}
	}

	logger.WithField("duration", time.Since(now)/time.Microsecond).Debug("Update")
	return result, nil
}

func (s *Session) Delete(ctx context.Context, id string, condition *request.Condition) (bool, error) {
	ctx, logger := log.ContextAndLoggerWithFields(ctx, log.Fields{"id": id, "condition": condition})

	if ctx == nil {
		return false, errors.New("context is missing")
	}
	if id == "" {
		return false, errors.New("id is missing")
	} else if !blob.IsValidID(id) {
		return false, errors.New("id is invalid")
	}
	if condition == nil {
		condition = request.NewCondition()
	} else if err := structureValidator.New().Validate(condition); err != nil {
		return false, errors.Wrap(err, "condition is invalid")
	}

	if s.IsClosed() {
		return false, errors.New("session closed")
	}

	now := time.Now()

	query := bson.M{
		"id": id,
	}
	if condition.Revision != nil {
		query["revision"] = *condition.Revision
	}
	set := bson.M{
		"modifiedTime": now.Truncate(time.Millisecond),
		"deletedTime":  now.Truncate(time.Millisecond),
	}
	unset := bson.M{}
	changeInfo, err := s.C().UpdateAll(query, s.ConstructUpdate(set, unset))
	if err != nil {
		logger.WithError(err).Error("Unable to delete blob")
		return false, errors.Wrap(err, "unable to delete blob")
	}

	logger.WithFields(log.Fields{"changeInfo": changeInfo, "duration": time.Since(now) / time.Microsecond}).Debug("Delete")
	return changeInfo.Updated > 0, nil
}

func (s *Session) Destroy(ctx context.Context, id string, condition *request.Condition) (bool, error) {
	ctx, logger := log.ContextAndLoggerWithFields(ctx, log.Fields{"id": id, "condition": condition})

	if ctx == nil {
		return false, errors.New("context is missing")
	}
	if id == "" {
		return false, errors.New("id is missing")
	} else if !blob.IsValidID(id) {
		return false, errors.New("id is invalid")
	}
	if condition == nil {
		condition = request.NewCondition()
	} else if err := structureValidator.New().Validate(condition); err != nil {
		return false, errors.Wrap(err, "condition is invalid")
	}

	if s.IsClosed() {
		return false, errors.New("session closed")
	}

	now := time.Now()

	query := bson.M{
		"id": id,
	}
	if condition.Revision != nil {
		query["revision"] = *condition.Revision
	}
	changeInfo, err := s.C().RemoveAll(query)
	if err != nil {
		logger.WithError(err).Error("Unable to destroy blob")
		return false, errors.Wrap(err, "unable to destroy blob")
	}

	logger.WithFields(log.Fields{"changeInfo": changeInfo, "duration": time.Since(now) / time.Microsecond}).Debug("Destroy")
	return changeInfo.Removed > 0, nil
}

func (s *Session) get(logger log.Logger, id string, condition *request.Condition, queryModifiers ...storeStructuredMongo.QueryModifier) (*blob.Blob, error) {
	logger = logger.WithFields(log.Fields{"id": id, "condition": condition})

	results := blob.BlobArray{}
	query := bson.M{
		"id": id,
	}
	if condition != nil && condition.Revision != nil {
		query["revision"] = *condition.Revision
	}
	query = storeStructuredMongo.ModifyQuery(query, queryModifiers...)
	err := s.C().Find(query).Limit(2).All(&results)
	if err != nil {
		logger.WithError(err).Error("Unable to get blob")
		return nil, errors.Wrap(err, "unable to get blob")
	}

	var result *blob.Blob
	switch len(results) {
	case 0:
		return nil, nil
	case 1:
		result = results[0]
	default:
		logger.Error("Multiple blobs found")
		result = results[0]
	}

	if result.Revision == nil {
		result.Revision = pointer.FromInt(0)
	}

	return result, nil
}
