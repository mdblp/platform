package log_test

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/tidepool-org/platform/errors"
	"github.com/tidepool-org/platform/log"
)

type Serializer struct {
	SerializeInvocations int
	SerializeInputs      []log.Fields
	SerializeOutputs     []error
}

func NewSerializer() *Serializer {
	return &Serializer{}
}

func (s *Serializer) Serialize(fields log.Fields) error {
	s.SerializeInvocations++

	s.SerializeInputs = append(s.SerializeInputs, fields)

	if len(s.SerializeOutputs) == 0 {
		panic("Unexpected invocation of Serialize on Serializer")
	}

	output := s.SerializeOutputs[0]
	s.SerializeOutputs = s.SerializeOutputs[1:]
	return output
}

func (s *Serializer) UnusedOutputsCount() int {
	return len(s.SerializeOutputs)
}

var _ = Describe("Logger", func() {
	var serializer *Serializer

	BeforeEach(func() {
		serializer = NewSerializer()
		Expect(serializer).ToNot(BeNil())
	})

	AfterEach(func() {
		Expect(serializer.UnusedOutputsCount()).To(Equal(0))
	})

	Context("NewLogger", func() {
		It("returns an error if the serializer is missing", func() {
			logger, err := log.NewLogger(nil, log.DefaultLevelRanks(), log.DefaultLevel())
			Expect(err).To(MatchError("serializer is missing"))
			Expect(logger).To(BeNil())
		})

		It("returns an error if the level ranks is missing", func() {
			logger, err := log.NewLogger(serializer, nil, log.DefaultLevel())
			Expect(err).To(MatchError("level ranks is missing"))
			Expect(logger).To(BeNil())
		})

		It("returns an error if the level is not found", func() {
			logger, err := log.NewLogger(serializer, log.DefaultLevelRanks(), log.Level("unknown"))
			Expect(err).To(MatchError("level not found"))
			Expect(logger).To(BeNil())
		})

		It("returns successfully", func() {
			Expect(log.NewLogger(serializer, log.DefaultLevelRanks(), log.DefaultLevel())).ToNot(BeNil())
		})
	})

	Context("with new logger", func() {
		var logger log.Logger

		BeforeEach(func() {
			var err error
			logger, err = log.NewLogger(serializer, log.DefaultLevelRanks(), log.DefaultLevel())
			Expect(err).ToNot(HaveOccurred())
			Expect(logger).ToNot(BeNil())
		})

		Context("Log", func() {
			It("does not invoke serializer if the level is unknown", func() {
				logger.Log(log.Level("unknown"), "Unknown Level Message")
			})

			It("does not invoke serializer if the level is not logging", func() {
				logger.Log(log.DebugLevel, "Not Logging Message")
			})

			Context("with disabled standard error", func() {
				var newFile *os.File
				var oldFile *os.File

				BeforeEach(func() {
					var err error
					newFile, err = ioutil.TempFile("", "")
					Expect(err).ToNot(HaveOccurred())
					Expect(newFile).ToNot(BeNil())
					oldFile = os.Stderr
					os.Stderr = newFile
				})

				AfterEach(func() {
					os.Stderr = oldFile
					os.Remove(newFile.Name())
				})

				It("fails silently if the serializer returns an error", func() {
					serializer.SerializeOutputs = []error{fmt.Errorf("test error")}
					logger.Log(log.WarnLevel, "Serializer Error Message")
					Expect(serializer.SerializeInputs).ToNot(BeEmpty())
				})
			})

			It("includes the expected fields in the expected format", func() {
				serializer.SerializeOutputs = []error{nil}
				logger.Log(log.WarnLevel, "Expected Message")
				Expect(serializer.SerializeInputs).To(HaveLen(1))
				serializeInput := serializer.SerializeInputs[0]
				Expect(serializeInput).To(HaveKey("dbl_caller"))
				Expect(serializeInput).To(HaveKeyWithValue("dbl_level", log.WarnLevel))
				Expect(serializeInput).To(HaveKey("dbl_time"))
				Expect(serializeInput).To(HaveKeyWithValue("dbl_message", "Expected Message"))
				serializedTime, ok := serializeInput["dbl_time"].(string)
				Expect(ok).To(BeTrue())
				loc, _ := time.LoadLocation("Local")
				parsedTime, err := time.ParseInLocation("2006-01-02T15:04:05.999999999", serializedTime, loc)
				Expect(err).ToNot(HaveOccurred())
				Expect(parsedTime).To(BeTemporally("~", time.Now(), time.Second))
				serializedCaller, ok := serializeInput["dbl_caller"].(*errors.Caller)
				Expect(ok).To(BeTrue())
				Expect(serializedCaller.Line).To(BeNumerically(">", 0))
				Expect(strings.HasSuffix(serializedCaller.File, "log/logger_test.go")).To(BeTrue())
			})

			It("does not include the message is it is an empty string", func() {
				serializer.SerializeOutputs = []error{nil}
				logger.Log(log.WarnLevel, "")
				Expect(serializer.SerializeInputs).To(HaveLen(1))
				Expect(serializer.SerializeInputs[0]).ToNot(HaveKey("dbl_message"))
			})
		})

		Context("with successful serialize and debug level", func() {
			BeforeEach(func() {
				serializer.SerializeOutputs = []error{nil}
				logger = logger.WithLevel(log.DebugLevel)
				Expect(logger).ToNot(BeNil())
			})

			Context("Debug", func() {
				It("logs with the expected level and message", func() {
					logger.Debug("Amazonian")
					Expect(serializer.SerializeInputs).To(HaveLen(1))
					Expect(serializer.SerializeInputs[0]).To(HaveKeyWithValue("dbl_level", log.DebugLevel))
					Expect(serializer.SerializeInputs[0]).To(HaveKeyWithValue("dbl_message", "Amazonian"))
				})
			})

			Context("Info", func() {
				It("logs with the expected level and message", func() {
					logger.Info("Bostonian")
					Expect(serializer.SerializeInputs).To(HaveLen(1))
					Expect(serializer.SerializeInputs[0]).To(HaveKeyWithValue("dbl_level", log.InfoLevel))
					Expect(serializer.SerializeInputs[0]).To(HaveKeyWithValue("dbl_message", "Bostonian"))
				})
			})

			Context("Warn", func() {
				It("logs with the expected level and message", func() {
					logger.Warn("Canadian")
					Expect(serializer.SerializeInputs).To(HaveLen(1))
					Expect(serializer.SerializeInputs[0]).To(HaveKeyWithValue("dbl_level", log.WarnLevel))
					Expect(serializer.SerializeInputs[0]).To(HaveKeyWithValue("dbl_message", "Canadian"))
				})
			})

			Context("Error", func() {
				It("logs with the expected level and message", func() {
					logger.Error("Dutch")
					Expect(serializer.SerializeInputs).To(HaveLen(1))
					Expect(serializer.SerializeInputs[0]).To(HaveKeyWithValue("dbl_level", log.ErrorLevel))
					Expect(serializer.SerializeInputs[0]).To(HaveKeyWithValue("dbl_message", "Dutch"))
				})
			})

			Context("Debugf", func() {
				It("logs with the expected level and message", func() {
					logger.Debugf("Amazonian %s", "Warrior")
					Expect(serializer.SerializeInputs).To(HaveLen(1))
					Expect(serializer.SerializeInputs[0]).To(HaveKeyWithValue("dbl_level", log.DebugLevel))
					Expect(serializer.SerializeInputs[0]).To(HaveKeyWithValue("dbl_message", "Amazonian Warrior"))
				})
			})

			Context("Infof", func() {
				It("logs with the expected level and message", func() {
					logger.Infof("Bostonian %s", "Cabbie")
					Expect(serializer.SerializeInputs).To(HaveLen(1))
					Expect(serializer.SerializeInputs[0]).To(HaveKeyWithValue("dbl_level", log.InfoLevel))
					Expect(serializer.SerializeInputs[0]).To(HaveKeyWithValue("dbl_message", "Bostonian Cabbie"))
				})
			})

			Context("Warnf", func() {
				It("logs with the expected level and message", func() {
					logger.Warnf("Canadian %s", "Skater")
					Expect(serializer.SerializeInputs).To(HaveLen(1))
					Expect(serializer.SerializeInputs[0]).To(HaveKeyWithValue("dbl_level", log.WarnLevel))
					Expect(serializer.SerializeInputs[0]).To(HaveKeyWithValue("dbl_message", "Canadian Skater"))
				})
			})

			Context("Errorf", func() {
				It("logs with the expected level and message", func() {
					logger.Errorf("Dutch %s", "Brothers")
					Expect(serializer.SerializeInputs).To(HaveLen(1))
					Expect(serializer.SerializeInputs[0]).To(HaveKeyWithValue("dbl_level", log.ErrorLevel))
					Expect(serializer.SerializeInputs[0]).To(HaveKeyWithValue("dbl_message", "Dutch Brothers"))
				})
			})

			Context("WithError", func() {
				It("does not include the error field if the error is missing", func() {
					logger.WithError(nil).Warn("European")
					Expect(serializer.SerializeInputs).To(HaveLen(1))
					Expect(serializer.SerializeInputs[0]).ToNot(HaveKey("dbl_error"))
				})

				It("deletes the error field if the error is missing", func() {
					logger.WithError(fmt.Errorf("euro error")).WithError(nil).Warn("European")
					Expect(serializer.SerializeInputs).To(HaveLen(1))
					Expect(serializer.SerializeInputs[0]).ToNot(HaveKey("dbl_error"))
				})

				It("does include the error field if the error is not missing", func() {
					logger.WithError(fmt.Errorf("euro error")).Warn("European")
					Expect(serializer.SerializeInputs).To(HaveLen(1))
					Expect(serializer.SerializeInputs[0]).To(HaveKey("dbl_error"))
				})
			})

			Context("WithField", func() {
				It("does not include the field if the key is missing", func() {
					logger.WithField("", "fish").Warn("Finnish")
					Expect(serializer.SerializeInputs).To(HaveLen(1))
					Expect(serializer.SerializeInputs[0]).ToNot(HaveKey(""))
				})

				It("does not include the field if the value is missing", func() {
					logger.WithField("sword", nil).Warn("Finnish")
					Expect(serializer.SerializeInputs).To(HaveLen(1))
					Expect(serializer.SerializeInputs[0]).ToNot(HaveKey("dbl_sword"))
				})

				It("deletes the field if the value is missing", func() {
					logger.WithField("sword", "fish").WithField("sword", nil).Warn("Finnish")
					Expect(serializer.SerializeInputs).To(HaveLen(1))
					Expect(serializer.SerializeInputs[0]).ToNot(HaveKey("dbl_sword"))
				})

				It("does include the field if the key and value are not missing", func() {
					logger.WithField("sword", "fish").Warn("Finnish")
					Expect(serializer.SerializeInputs).To(HaveLen(1))
					Expect(serializer.SerializeInputs[0]).To(HaveKeyWithValue("dbl_sword", "fish"))
				})
			})

			Context("WithFields", func() {
				It("does include the field if the key and value are not missing", func() {
					logger.WithFields(log.Fields{"": "Nein", "nope": nil, "yep": "Ja"}).Warn("German")
					Expect(serializer.SerializeInputs).To(HaveLen(1))
					Expect(serializer.SerializeInputs[0]).ToNot(HaveKey(""))
					Expect(serializer.SerializeInputs[0]).ToNot(HaveKey("dbl_nope"))
					Expect(serializer.SerializeInputs[0]).To(HaveKeyWithValue("dbl_yep", "Ja"))
				})

				It("deletes the field if the value is missing", func() {
					logger.WithFields(log.Fields{"nope": "Nein"}).WithFields(log.Fields{"nope": nil}).Warn("German")
					Expect(serializer.SerializeInputs).To(HaveLen(1))
					Expect(serializer.SerializeInputs[0]).ToNot(HaveKey("dbl_nope"))
				})
			})

			Context("WithLevelRank", func() {
				It("adds the specified level and rank", func() {
					level := log.Level("new")
					logger = logger.WithLevelRank(level, 90).WithLevel(level)
					Expect(logger).ToNot(BeNil())
					logger.Debug("Should Not Serialize")
					logger.Log(level, "WithLevelRank Message")
					Expect(serializer.SerializeInputs).To(HaveLen(1))
					Expect(serializer.SerializeInputs[0]).To(HaveKeyWithValue("dbl_level", level))
					Expect(serializer.SerializeInputs[0]).To(HaveKeyWithValue("dbl_message", "WithLevelRank Message"))
				})
			})

			Context("WithLevelRanks", func() {
				It("adds the specified level ranks", func() {
					level := log.Level("new")
					logger = logger.WithLevelRanks(log.LevelRanks{level: 30, log.Level("other"): 0}).WithLevel(level)
					Expect(logger).ToNot(BeNil())
					logger.Debug("Should Not Serialize")
					logger.Log(level, "WithLevelRanks Message")
					Expect(serializer.SerializeInputs).To(HaveLen(1))
					Expect(serializer.SerializeInputs[0]).To(HaveKeyWithValue("dbl_level", level))
					Expect(serializer.SerializeInputs[0]).To(HaveKeyWithValue("dbl_message", "WithLevelRanks Message"))
				})
			})

			Context("WithLevel", func() {
				It("uses the current level if the specified level is unknown", func() {
					logger = logger.WithLevel(log.Level("unknown"))
					Expect(logger).ToNot(BeNil())
					Expect(logger.Level()).To(Equal(log.DebugLevel))
					logger.Debug("WithLevel Message")
					Expect(serializer.SerializeInputs).To(HaveLen(1))
					Expect(serializer.SerializeInputs[0]).To(HaveKeyWithValue("dbl_level", log.DebugLevel))
					Expect(serializer.SerializeInputs[0]).To(HaveKeyWithValue("dbl_message", "WithLevel Message"))
				})

				It("adds the specified level", func() {
					logger = logger.WithLevel(log.InfoLevel)
					Expect(logger).ToNot(BeNil())
					logger.Debug("Should Not Serialize")
					logger.Warn("WithLevel Message")
					Expect(serializer.SerializeInputs).To(HaveLen(1))
					Expect(serializer.SerializeInputs[0]).To(HaveKeyWithValue("dbl_level", log.WarnLevel))
					Expect(serializer.SerializeInputs[0]).To(HaveKeyWithValue("dbl_message", "WithLevel Message"))
				})
			})
		})

		Context("Level", func() {
			It("returns the current level", func() {
				Expect(logger.Level()).To(Equal(log.WarnLevel))
			})

			It("returns the level after being set", func() {
				logger = logger.WithLevel(log.DebugLevel)
				Expect(logger).ToNot(BeNil())
				Expect(logger.Level()).To(Equal(log.DebugLevel))
			})

			It("returns the level after a new level is added and set", func() {
				level := log.Level("new")
				logger = logger.WithLevelRank(level, 55).WithLevel(level)
				Expect(logger).ToNot(BeNil())
				Expect(logger.Level()).To(Equal(level))
			})
		})
	})
})
