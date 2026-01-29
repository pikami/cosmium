package continuationtoken

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/pikami/cosmium/internal/logger"
)

type ContinuationTokenExternal struct {
	Token string `json:"token"`
	Range struct {
		Min string `json:"min"`
		Max string `json:"max"`
	} `json:"range"`
}

type ContinuationToken struct {
	Token struct {
		ResourceId   string // RID
		PageIndex    int    // RT
		TotalResults int    // TRC
		ISV          int    // ISV
		IEO          int    // IEO
		QCF          int    // QCF
		LR           int    // LR
	}
	Range struct {
		Min string
		Max string
	}
}

func Generate(resourceid string, pageIndex int, totalResults int) ContinuationToken {
	ct := ContinuationToken{}
	ct.Token.ResourceId = resourceid
	ct.Token.PageIndex = pageIndex
	ct.Token.TotalResults = totalResults
	ct.Token.ISV = 2
	ct.Token.IEO = 65567
	ct.Token.QCF = 8
	ct.Token.LR = 1
	ct.Range.Min = ""
	ct.Range.Max = "FF"

	return ct
}

func GenerateDefault(resourceid string) ContinuationToken {
	return Generate(resourceid, 0, 0)
}

func (ct *ContinuationToken) ToString() string {
	token := fmt.Sprintf(
		"-RID:~%s#RT:%d#TRC:%d#ISV:%d#IEO:%d#QCF:%d#LR:%d",
		ct.Token.ResourceId,
		ct.Token.PageIndex,
		ct.Token.TotalResults,
		ct.Token.ISV,
		ct.Token.IEO,
		ct.Token.QCF,
		ct.Token.LR,
	)

	ect := ContinuationTokenExternal{}
	ect.Token = token
	ect.Range.Min = ct.Range.Min
	ect.Range.Max = ct.Range.Max

	json, err := json.Marshal(ect)
	if err != nil {
		logger.Error(err, "failed to marshal continuation token")
		return ""
	}

	return string(json)
}

func FromString(token string) ContinuationToken {
	ect := ContinuationTokenExternal{}
	err := json.Unmarshal([]byte(token), &ect)
	if err != nil {
		logger.Error(err, "failed to unmarshal continuation token")
		return ContinuationToken{}
	}

	ct, err := parseContinuationToken(ect.Token, ect.Range.Min, ect.Range.Max)
	if err != nil {
		logger.Error(err, "failed to parse continuation token")
		return ContinuationToken{}
	}

	return *ct
}

func parseContinuationToken(token string, minRange string, maxRange string) (*ContinuationToken, error) {
	const prefix = "-RID:~"
	if !strings.HasPrefix(token, prefix) {
		return nil, fmt.Errorf("invalid token prefix")
	}

	parts := strings.Split(token[len(prefix):], "#")
	if len(parts) != 7 {
		return nil, fmt.Errorf("invalid token format: expected 7 fields, got %d", len(parts))
	}

	ct := &ContinuationToken{}

	ct.Token.ResourceId = parts[0]

	parseIntField := func(part, key string) (int, error) {
		if !strings.HasPrefix(part, key+":") {
			return 0, fmt.Errorf("expected %s field", key)
		}
		return strconv.Atoi(strings.TrimPrefix(part, key+":"))
	}

	var err error

	if ct.Token.PageIndex, err = parseIntField(parts[1], "RT"); err != nil {
		return nil, err
	}
	if ct.Token.TotalResults, err = parseIntField(parts[2], "TRC"); err != nil {
		return nil, err
	}
	if ct.Token.ISV, err = parseIntField(parts[3], "ISV"); err != nil {
		return nil, err
	}
	if ct.Token.IEO, err = parseIntField(parts[4], "IEO"); err != nil {
		return nil, err
	}
	if ct.Token.QCF, err = parseIntField(parts[5], "QCF"); err != nil {
		return nil, err
	}
	if ct.Token.LR, err = parseIntField(parts[6], "LR"); err != nil {
		return nil, err
	}

	ct.Range.Min = minRange
	ct.Range.Max = maxRange

	return ct, nil
}
