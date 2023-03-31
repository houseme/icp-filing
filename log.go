/*
 *  Copyright icp-filing Author(https://houseme.github.io/icp-filing/). All Rights Reserved.
 *
 *  Licensed under the Apache License, Version 2.0 (the "License");
 *  you may not use this file except in compliance with the License.
 *  You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 *  Unless required by applicable law or agreed to in writing, software
 *  distributed under the License is distributed on an "AS IS" BASIS,
 *  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *  See the License for the specific language governing permissions and
 *  limitations under the License.
 *
 *  You can obtain one at https://github.com/houseme/icp-filing.
 */

package filling

import (
	"os"
	"time"

	"github.com/cloudwego/hertz/pkg/common/hlog"
	hertzzap "github.com/hertz-contrib/logger/zap"
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func initLog(path string, level hlog.Level) {
	dynamicLevel := zap.NewAtomicLevel()
	dynamicLevel.SetLevel(zap.DebugLevel)
	logger := hertzzap.NewLogger(
		hertzzap.WithCores([]hertzzap.CoreConfig{
			{
				Enc: zapcore.NewConsoleEncoder(humanEncoderConfig()),
				Ws:  os.Stdout,
				Lvl: dynamicLevel,
			},
			{
				Enc: zapcore.NewJSONEncoder(humanEncoderConfig()),
				Ws:  getWriteSyncer(path + "/all.log"),
				Lvl: zap.NewAtomicLevelAt(zapcore.DebugLevel),
			},
			{
				Enc: zapcore.NewJSONEncoder(humanEncoderConfig()),
				Ws:  getWriteSyncer(path + "/debug.log"),
				Lvl: zap.NewAtomicLevelAt(zapcore.LevelOf(
					zap.LevelEnablerFunc(func(lev zapcore.Level) bool {
						return lev == zap.DebugLevel
					}))),
			},
			{
				Enc: zapcore.NewJSONEncoder(humanEncoderConfig()),
				Ws:  getWriteSyncer(path + "/info.log"),
				Lvl: zap.NewAtomicLevelAt(zapcore.LevelOf(
					zap.LevelEnablerFunc(func(lev zapcore.Level) bool {
						return lev == zap.InfoLevel
					}))),
			},
			{
				Enc: zapcore.NewJSONEncoder(humanEncoderConfig()),
				Ws:  getWriteSyncer(path + "/warn.log"),
				Lvl: zap.NewAtomicLevelAt(zapcore.LevelOf(
					zap.LevelEnablerFunc(func(lev zapcore.Level) bool {
						return lev == zap.WarnLevel
					}))),
			},
			{
				Enc: zapcore.NewJSONEncoder(humanEncoderConfig()),
				Ws:  getWriteSyncer(path + "/error.log"),
				Lvl: zap.NewAtomicLevelAt(zapcore.LevelOf(
					zap.LevelEnablerFunc(func(lev zapcore.Level) bool {
						return lev >= zap.ErrorLevel
					}))),
			},
		}...),
	)
	defer logger.Sync()
	hlog.SetLogger(logger)
	hlog.SetLevel(level)
	hlog.Infof("filing start %s", time.Now().String())
}

// humanEncoderConfig copy from zap
func humanEncoderConfig() zapcore.EncoderConfig {
	cfg := testEncoderConfig()
	cfg.EncodeTime = zapcore.ISO8601TimeEncoder
	cfg.EncodeLevel = zapcore.CapitalLevelEncoder
	cfg.EncodeDuration = zapcore.StringDurationEncoder
	return cfg
}

func getWriteSyncer(file string) zapcore.WriteSyncer {
	lumberJackLogger := &lumberjack.Logger{
		Filename:   file,
		MaxSize:    10,
		MaxBackups: 50000,
		MaxAge:     1000,
		Compress:   true,
		LocalTime:  true,
	}
	return zapcore.AddSync(lumberJackLogger)
}

// testEncoderConfig encoder config for testing, copy from zap
func testEncoderConfig() zapcore.EncoderConfig {
	return zapcore.EncoderConfig{
		MessageKey:     "msg",
		LevelKey:       "level",
		NameKey:        "name",
		TimeKey:        "ts",
		CallerKey:      "caller",
		FunctionKey:    "func",
		StacktraceKey:  "stacktrace",
		LineEnding:     "\n",
		EncodeTime:     zapcore.EpochTimeEncoder,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
}
