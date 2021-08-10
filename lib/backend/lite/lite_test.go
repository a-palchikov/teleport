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
	"os"
	"testing"
	"time"

	"github.com/gravitational/teleport/lib/backend"
	"github.com/gravitational/teleport/lib/backend/test"
	"github.com/gravitational/teleport/lib/fixtures"
	"github.com/gravitational/teleport/lib/utils"

	"github.com/jonboulle/clockwork"
	"gopkg.in/check.v1"
)

func TestMain(m *testing.M) {
	utils.InitLoggerForTests()
	os.Exit(m.Run())
}

func TestLite(t *testing.T) { check.TestingT(t) }

type LiteSuite struct {
	bk    *Backend
	suite test.BackendSuite
}

var _ = check.Suite(&LiteSuite{})

func (s *LiteSuite) SetUpTest(c *check.C) {
	var err error
	clock := clockwork.NewFakeClock()
	s.bk, err = NewWithConfig(context.Background(), Config{
		Path:             c.MkDir(),
		PollStreamPeriod: 300 * time.Millisecond,
		Clock:            clock,
	})
	c.Assert(err, check.IsNil)
	s.suite.B = s.bk
	s.suite.Clock = clock
}

func (s *LiteSuite) TearDownTest(c *check.C) {
	if s.bk != nil {
		c.Assert(s.bk.Close(), check.IsNil)
	}
}

func (s *LiteSuite) TestCRUD(c *check.C) {
	s.suite.CRUD(c)
}

func (s *LiteSuite) TestRange(c *check.C) {
	s.suite.Range(c)
}

func (s *LiteSuite) TestCompareAndSwap(c *check.C) {
	s.suite.CompareAndSwap(c)
}

func (s *LiteSuite) TestExpiration(c *check.C) {
	s.suite.Expiration(c)
}

func (s *LiteSuite) TestKeepAlive(c *check.C) {
	s.suite.KeepAlive(c)
}

func (s *LiteSuite) TestEvents(c *check.C) {
	s.suite.Events(c)
}

func (s *LiteSuite) TestWatchersClose(c *check.C) {
	s.suite.WatchersClose(c)
}

func (s *LiteSuite) TestDeleteRange(c *check.C) {
	s.suite.DeleteRange(c)
}

func (s *LiteSuite) TestPutRange(c *check.C) {
	s.suite.PutRange(c)
}

func (s *LiteSuite) TestLocking(c *check.C) {
	s.suite.Locking(c, s.bk)
}

// Import tests importing values
func (s *LiteSuite) TestImport(c *check.C) {
	ctx := context.Background()
	prefix := test.MakePrefix()

	imported, err := s.bk.Imported(ctx)
	c.Assert(err, check.IsNil)
	c.Assert(imported, check.Equals, false)

	// add one element that should not show up
	items := []backend.Item{
		{Key: prefix("/prefix/a"), Value: []byte("val a")},
		{Key: prefix("/prefix/b"), Value: []byte("val b")},
		{Key: prefix("/prefix/a"), Value: []byte("val a")},
	}
	err = s.bk.Import(ctx, items)
	c.Assert(err, check.IsNil)

	// prefix range fetch
	result, err := s.bk.GetRange(ctx, prefix("/prefix"), backend.RangeEnd(prefix("/prefix")), backend.NoLimit)
	c.Assert(err, check.IsNil)
	expected := []backend.Item{
		{Key: prefix("/prefix/a"), Value: []byte("val a")},
		{Key: prefix("/prefix/b"), Value: []byte("val b")},
	}
	test.ExpectItems(c, result.Items, expected)

	imported, err = s.bk.Imported(ctx)
	c.Assert(err, check.IsNil)
	c.Assert(imported, check.Equals, true)

	err = s.bk.Import(ctx, items)
	fixtures.ExpectAlreadyExists(c, err)

	imported, err = s.bk.Imported(ctx)
	c.Assert(err, check.IsNil)
	c.Assert(imported, check.Equals, true)
}
