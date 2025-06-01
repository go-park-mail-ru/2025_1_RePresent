package attempt

import (
	"errors"
	"testing"
	"time"

	"github.com/go-redis/redismock/v9"
	"github.com/stretchr/testify/assert"
)

func Test_ResetAttemptsByUserID_Success(t *testing.T) {
	db, mock := redismock.NewClientMock()
	repo := &AttemptRepository{Client: db, Ttl: time.Minute}

	mock.ExpectDel("attempts:10").SetVal(1)
	err := repo.ResetAttemptsByUserID(10)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func Test_ResetAttemptsByUserID_Error(t *testing.T) {
	db, mock := redismock.NewClientMock()
	repo := &AttemptRepository{Client: db, Ttl: time.Minute}

	mock.ExpectDel("attempts:5").SetErr(errors.New("del error"))
	err := repo.ResetAttemptsByUserID(5)
	assert.Error(t, err)
}

func Test_IncrementAttemptsByUserID_FirstAndThenNoTTL(t *testing.T) {
	db, mock := redismock.NewClientMock()
	repo := &AttemptRepository{Client: db, Ttl: time.Minute}

	mock.ExpectIncr("attempts:1").SetVal(1)
	mock.ExpectExpire("attempts:1", repo.Ttl).SetVal(true)
	err := repo.IncrementAttemptsByUserID(1)
	assert.NoError(t, err)

	mock.ExpectIncr("attempts:1").SetVal(2)
	err = repo.IncrementAttemptsByUserID(1)
	assert.NoError(t, err)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func Test_GetAttemptsByUserID_ZeroAndValue(t *testing.T) {
	db, mock := redismock.NewClientMock()
	repo := &AttemptRepository{Client: db, Ttl: time.Minute}

	mock.ExpectGet("attempts:2").RedisNil()
	cnt, err := repo.GetAttemptsByUserID(2)
	assert.NoError(t, err)
	assert.Equal(t, 0, cnt)

	mock.ExpectGet("attempts:2").SetVal("abc")
	_, err = repo.GetAttemptsByUserID(2)
	assert.Error(t, err)

	mock.ExpectGet("attempts:3").SetVal("7")
	cnt, err = repo.GetAttemptsByUserID(3)
	assert.NoError(t, err)
	assert.Equal(t, 7, cnt)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func Test_DecrementAttemptsByUserID_EvalAndError(t *testing.T) {
	db, mock := redismock.NewClientMock()
	repo := &AttemptRepository{Client: db, Ttl: time.Minute}

	mock.ExpectEval(redisScript, []string{"attempts:9"}).SetVal(int64(0))
	err := repo.DecrementAttemptsByUserID(9)

	mock.ExpectEval(redisScript, []string{"attempts:9"}).SetErr(errors.New("eval error"))
	err = repo.DecrementAttemptsByUserID(9)
	assert.Error(t, err)
}

const redisScript = `
        local current = redis.call("DECR", KEYS[1])
        
        -- delete if value <= 0
        if current <= 0 then
            redis.call("DEL", KEYS[1])
            return 0
        end
        
        return current
    `
