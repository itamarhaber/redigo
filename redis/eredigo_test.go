package redis_test

import (
	"fmt"
	"testing"

	"github.com/gomodule/redigo/redis"
)

func assertEqual(t *testing.T, a interface{}, b interface{}) {
	if a == b {
		return
	}
	t.Fatal(fmt.Sprintf("assertEqual failed: %v != %v", a, b))
}

func TestEredigoDoPing(t *testing.T) {
	c, err := redis.Dial("eredis", "")
	if err != nil {
		t.Errorf("Could not dial: %s", err.Error())
		t.Fail()
	}
	defer c.Close()

	rep, err := redis.String(c.Do("PING"))
	if err != nil {
		t.Errorf("Could not PING: %s", err.Error())
		t.Fail()
	}

	assertEqual(t, rep, "PONG")
}

func TestEredigoDoTwoCommands(t *testing.T) {
	c, err := redis.Dial("eredis", "")
	if err != nil {
		t.Errorf("Could not dial: %s", err.Error())
		t.Fail()
	}
	defer c.Close()

	rep, err := redis.String(c.Do("PING"))
	if err != nil {
		t.Errorf("Could not PING once: %s", err.Error())
		t.Fail()
	}

	assertEqual(t, rep, "PONG")

	rep, err = redis.String(c.Do("PING", "TESTUNG"))
	if err != nil {
		t.Errorf("Could not PING twice: %s", err.Error())
		t.Fail()
	}

	assertEqual(t, rep, "TESTUNG")
}

func TestEredigoPipeline(t *testing.T) {
	c, err := redis.Dial("eredis", "")
	if err != nil {
		t.Errorf("Could not dial: %s", err.Error())
		t.Fail()
	}
	defer c.Close()

	err = c.Send("ECHO", "foo")
	if err != nil {
		t.Errorf("Could not ECHO foo: %s", err.Error())
		t.Fail()
	}

	err = c.Send("ECHO", "bar")
	if err != nil {
		t.Errorf("Could not ECHO bar: %s", err.Error())
		t.Fail()
	}

	rep, err := redis.Strings(c.Do(""))
	if err != nil {
		t.Errorf("Could not do echos pipeline: %s", err.Error())
		t.Fail()
	}

	assertEqual(t, len(rep), 2)
	assertEqual(t, rep[0], "foo")
	assertEqual(t, rep[1], "bar")
}

func TestEredigoTransaction(t *testing.T) {
	c, err := redis.Dial("eredis", "")
	if err != nil {
		t.Errorf("Could not dial: %s", err.Error())
		t.Fail()
	}
	defer c.Close()

	err = c.Send("MULTI")
	if err != nil {
		t.Errorf("Could not MULTI: %s", err.Error())
		t.Fail()
	}

	err = c.Send("FLUSHALL")
	if err != nil {
		t.Errorf("Could not FLUSHALL: %s", err.Error())
		t.Fail()
	}

	err = c.Send("SET", "foo", "bar")
	if err != nil {
		t.Errorf("Could not SET foo to bar: %s", err.Error())
		t.Fail()
	}

	err = c.Send("GET", "foo")
	if err != nil {
		t.Errorf("Could not GET foo: %s", err.Error())
		t.Fail()
	}

	rep, err := redis.Strings(c.Do("EXEC"))
	if err != nil {
		t.Errorf("Could not EXEC: %s", err.Error())
		t.Fail()
	}

	assertEqual(t, len(rep), 3)
	assertEqual(t, rep[0], "OK")
	assertEqual(t, rep[1], "OK")
	assertEqual(t, rep[2], "bar")
}

func TestEredigoScript(t *testing.T) {
	c, err := redis.Dial("eredis", "")
	if err != nil {
		t.Errorf("Could not dial: %s", err.Error())
		t.Fail()
	}
	defer c.Close()

	s := redis.NewScript(0, `
		local rep = {}
		for _, v in ipairs(ARGV) do
			table.insert(rep, 1, v)
		end
		return rep
		`)

	rep, err := redis.Strings(s.Do(c, "foo", "bar"))
	if err != nil {
		t.Errorf("Could not EVAL: %s", err.Error())
		t.Fail()
	}

	assertEqual(t, len(rep), 2)
	assertEqual(t, rep[0], "bar")
	assertEqual(t, rep[1], "foo")

}
