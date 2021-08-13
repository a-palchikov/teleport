/*
Copyright 2018-2019 Gravitational, Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package lite

import (
	"context"
	"time"

	"github.com/gravitational/teleport/lib/backend/test"

	"github.com/jonboulle/clockwork"
	"gopkg.in/check.v1"
)

type LiteMemSuite struct {
	bk    *Backend
	clock clockwork.FakeClock
	suite test.BackendSuite
}

var _ = check.Suite(&LiteMemSuite{})

func (s *LiteMemSuite) SetUpTest(c *check.C) {
	s.clock = clockwork.NewFakeClock()
	s.bk = s.newBackend(c)
	s.suite.B = s.bk
	s.suite.Clock = s.clock
}

func (s *LiteMemSuite) newBackend(c *check.C) *Backend {
	m, err := NewWithConfig(context.Background(), Config{
		Memory:           true,
		PollStreamPeriod: 300 * time.Millisecond,
		Clock:            s.clock,
	})
	c.Assert(err, check.IsNil)
	return m
}

func (s *LiteMemSuite) TearDownTest(c *check.C) {
	if s.bk != nil {
		c.Assert(s.bk.Close(), check.IsNil)
	}
}

func (s *LiteMemSuite) TestCRUD(c *check.C) {
	s.suite.CRUD(c)
}

func (s *LiteMemSuite) TestRange(c *check.C) {
	s.suite.Range(c)
}

func (s *LiteMemSuite) TestCompareAndSwap(c *check.C) {
	s.suite.CompareAndSwap(c)
}

func (s *LiteMemSuite) TestExpiration(c *check.C) {
	s.suite.Expiration(c)
}

func (s *LiteMemSuite) TestKeepAlive(c *check.C) {
	s.suite.KeepAlive(c)
}

func (s *LiteMemSuite) TestEvents(c *check.C) {
	s.suite.Events(c)
}

func (s *LiteMemSuite) TestWatchersClose(c *check.C) {
	s.suite.WatchersClose(c)
}

func (s *LiteMemSuite) TestDeleteRange(c *check.C) {
	s.suite.DeleteRange(c)
}

func (s *LiteMemSuite) TestPutRange(c *check.C) {
	s.suite.PutRange(c)
}

func (s *LiteMemSuite) TestLocking(c *check.C) {
	s.suite.Locking(c, s.bk)
}

func (s *LiteMemSuite) TestConcurrentOperations(c *check.C) {
	bk := s.newBackend(c)
	defer bk.Close()
	s.suite.ConcurrentOperations(c, bk)
}

func (s *LiteMemSuite) TestMirror(c *check.C) {
	mem, err := NewWithConfig(context.Background(), Config{
		Memory:           true,
		Mirror:           true,
		PollStreamPeriod: 300 * time.Millisecond,
	})
	defer mem.Close()
	c.Assert(err, check.IsNil)
	s.suite.Mirror(c, mem)
}
