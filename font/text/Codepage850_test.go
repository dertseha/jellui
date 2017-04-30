package text

import (
	check "gopkg.in/check.v1"
)

type Codepage850Suite struct {
	cp Codepage
}

var _ = check.Suite(&Codepage850Suite{})

func (suite *Codepage850Suite) SetUpTest(c *check.C) {
	suite.cp = Codepage850()
}

func (suite *Codepage850Suite) TestEncode(c *check.C) {
	result := suite.cp.Encode("ä")

	c.Check(result, check.DeepEquals, []byte{132, 0x00})
}

func (suite *Codepage850Suite) TestDecode(c *check.C) {
	result := suite.cp.Decode([]byte{212, 225, 0x00})

	c.Check(result, check.Equals, "Èß")
}
