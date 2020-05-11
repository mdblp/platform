export TIDEPOOL_ENV="local"
export TIDEPOOL_LOGGER_LEVEL="debug"
export TIDEPOOL_STORE_ADDRESSES="localhost"
export TIDEPOOL_STORE_DATABASE="tidepool"
export TIDEPOOL_STORE_TLS="false"
export TIDEPOOL_SERVER_TLS="false"

export TIDEPOOL_CONFIRMATION_STORE_DATABASE="confirm"
export TIDEPOOL_DEPRECATED_DATA_STORE_DATABASE="data"
export TIDEPOOL_MESSAGE_STORE_DATABASE="messages"
export TIDEPOOL_PERMISSION_STORE_DATABASE="gatekeeper"
export TIDEPOOL_PERMISSION_STORE_SECRET="This secret is used to encrypt the groupId stored in the DB for gatekeeper"
export TIDEPOOL_PROFILE_STORE_DATABASE="seagull"
export TIDEPOOL_SESSION_STORE_DATABASE="user"
export TIDEPOOL_SYNC_TASK_STORE_DATABASE="data"
export TIDEPOOL_TASK_QUEUE_DELAY="5"
export TIDEPOOL_TASK_QUEUE_WORKERS="5"
export TIDEPOOL_USER_STORE_DATABASE="user"
export TIDEPOOL_USER_STORE_PASSWORD_SALT="ADihSEI7tOQQP9xfXMO9HfRpXKu1NpIJ"

export TIDEPOOL_AUTH_CLIENT_ADDRESS="http://localhost:8009"
export TIDEPOOL_BLOB_CLIENT_ADDRESS="http://localhost:8009"
export TIDEPOOL_DATA_CLIENT_ADDRESS="http://localhost:8009"
export TIDEPOOL_DATA_SOURCE_CLIENT_ADDRESS="http://localhost:8009"
export TIDEPOOL_IMAGE_CLIENT_ADDRESS="http://localhost:8009"
export TIDEPOOL_PERMISSION_CLIENT_ADDRESS="http://localhost:8009"
export TIDEPOOL_TASK_CLIENT_ADDRESS="http://localhost:8009"
export TIDEPOOL_USER_CLIENT_ADDRESS="http://localhost:8009"

export TIDEPOOL_AUTH_CLIENT_EXTERNAL_ADDRESS="http://localhost:8009"
export TIDEPOOL_AUTH_CLIENT_EXTERNAL_SERVER_SESSION_TOKEN_SECRET="This needs to be the same secret everywhere. YaHut75NsK1f9UKUXuWqxNN0RUwHFBCy"

export TIDEPOOL_AUTH_SERVICE_SERVER_ADDRESS=":9222"
export TIDEPOOL_BLOB_SERVICE_SERVER_ADDRESS=":9225"
export TIDEPOOL_DATA_SERVICE_SERVER_ADDRESS=":9220"
export TIDEPOOL_IMAGE_SERVICE_SERVER_ADDRESS=":9226"
export TIDEPOOL_NOTIFICATION_SERVICE_SERVER_ADDRESS=":9223"
export TIDEPOOL_TASK_SERVICE_SERVER_ADDRESS=":9224"
export TIDEPOOL_USER_SERVICE_SERVER_ADDRESS=":9221"

export TIDEPOOL_AUTH_SERVICE_DOMAIN="localhost"

export TIDEPOOL_BLOB_SERVICE_UNSTRUCTURED_STORE_TYPE="file"
export TIDEPOOL_BLOB_SERVICE_UNSTRUCTURED_STORE_FILE_DIRECTORY="_data/blobs"

export TIDEPOOL_IMAGE_SERVICE_UNSTRUCTURED_STORE_TYPE="file"
export TIDEPOOL_IMAGE_SERVICE_UNSTRUCTURED_STORE_FILE_DIRECTORY="_data/images"

export TIDEPOOL_AUTH_SERVICE_SECRET="Service secret used for interservice requests with the auth service"
export TIDEPOOL_BLOB_SERVICE_SECRET="Service secret used for interservice requests with the blob service"
export TIDEPOOL_DATA_SERVICE_SECRET="Service secret used for interservice requests with the data service"
export TIDEPOOL_IMAGE_SERVICE_SECRET="Service secret used for interservice requests with the image service"
export TIDEPOOL_NOTIFICATION_SERVICE_SECRET="Service secret used for interservice requests with the notification service"
export TIDEPOOL_TASK_SERVICE_SECRET="Service secret used for interservice requests with the task service"
export TIDEPOOL_USER_SERVICE_SECRET="Service secret used for interservice requests with the user service"
