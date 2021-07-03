package connectionQueue

import (
	"echoServer/models"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestConnectionQueue_Len(t *testing.T) {
	tests := []struct {
		name string
		cq   ConnectionQueue
		want int
	}{
		{
			name: "empty queue",
			cq:   NewConnectionQueue(),
			want: 0,
		},
		{
			name: "two element queue",
			cq:   NewConnectionQueue(),
			want: 2,
		},
	}
	for i := 0; i < 2; i++ {
		tests[1].cq = append(tests[1].cq, &models.Connection{
			Conn:       nil,
			LastUpdate: time.Time{},
			Index:      0,
		})
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.cq.Len()
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestConnectionQueue_Less(t *testing.T) {
	type args struct {
		i int
		j int
	}
	tests := []struct {
		name string
		cq   ConnectionQueue
		args args
		want bool
	}{
		{
			name: "bigger priority",
			cq:   NewConnectionQueue(),
			args: args{
				i: 0,
				j: 1,
			},
			want: false,
		},

		{
			name: "equal",
			cq:   NewConnectionQueue(),
			args: args{
				i: 0,
				j: 1,
			},
			want: false,
		},

		{
			name: "lower priority",
			cq:   NewConnectionQueue(),
			args: args{
				i: 0,
				j: 1,
			},
			want: true,
		},
	}

	tests[0].cq = append(tests[0].cq, &models.Connection{
		LastUpdate: time.Now().Add(-time.Second * 10),
	})
	tests[0].cq = append(tests[0].cq, &models.Connection{
		LastUpdate: time.Now(),
	})
	tests[1].cq = append(tests[1].cq, &models.Connection{
		LastUpdate: time.Date(2000, 0, 0, 0, 0, 0, 0, time.Now().Location()),
	})
	tests[1].cq = append(tests[1].cq, &models.Connection{
		LastUpdate: time.Date(2000, 0, 0, 0, 0, 0, 0, time.Now().Location()),
	})
	tests[2].cq = append(tests[2].cq, &models.Connection{
		LastUpdate: time.Date(2000, 0, 0, 0, 0, 0, 0, time.Now().Location()),
	})
	tests[2].cq = append(tests[2].cq, &models.Connection{
		LastUpdate: time.Date(1990, 0, 0, 0, 0, 0, 0, time.Now().Location()),
	})

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.cq.Less(tt.args.i, tt.args.j)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestConnectionQueue_Pop(t *testing.T) {
	tests := []struct {
		name string
		cq   ConnectionQueue
		want *models.Connection
	}{
		{
			name: "empty",
			cq:   NewConnectionQueue(),
			want: nil,
		},
		{
			name: "one element",
			cq:   NewConnectionQueue(),
			want: &models.Connection{
				Conn:       nil,
				LastUpdate: time.Date(2000, 0, 0, 0, 0, 0, 0, time.Now().Location()),
				Index:      -1,
			},
		},
		{
			name: "two element",
			cq:   NewConnectionQueue(),
			want: &models.Connection{
				Conn:       nil,
				LastUpdate: time.Date(1999, 0, 0, 0, 0, 0, 0, time.Now().Location()),
				Index:      -1,
			},
		},
	}

	tests[1].cq.Push(&models.Connection{
		LastUpdate: time.Date(2000, 0, 0, 0, 0, 0, 0, time.Now().Location()),
	})
	tests[2].cq.Push(&models.Connection{
		LastUpdate: time.Date(1999, 0, 0, 0, 0, 0, 0, time.Now().Location()),
	})
	tests[2].cq.Push(&models.Connection{
		LastUpdate: time.Date(2000, 0, 0, 0, 0, 0, 0, time.Now().Location()),
	})

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.cq.Pop()
			gotConverted, ok := got.(*models.Connection)
			if !ok {
				gotConverted = nil
			}
			assert.Equal(t, tt.want, gotConverted)
		})
	}
}

func TestConnectionQueue_Push(t *testing.T) {
	type args struct {
		x interface{}
	}
	tests := []struct {
		name string
		cq   ConnectionQueue
		args args
		want *models.Connection
	}{
		{
			name: "push one and get one",
			cq:   NewConnectionQueue(),
			args: args{x: &models.Connection{
				Conn:       nil,
				LastUpdate: time.Date(2000, 0, 0, 0, 0, 0, 0, time.Now().Location()),
				Index:      0}},
			want: &models.Connection{
				Conn:       nil,
				LastUpdate: time.Date(2000, 0, 0, 0, 0, 0, 0, time.Now().Location()),
				Index:      -1,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.cq.Push(tt.args.x)
			got := tt.cq.Pop()
			assert.Equal(t, tt.want, got.(*models.Connection))
		})
	}
}

func TestConnectionQueue_Update(t *testing.T) {
	type args struct {
		item       *models.Connection
		lastUpdate time.Time
	}
	tests := []struct {
		name string
		cq   ConnectionQueue
		args args
		want *models.Connection
	}{
		{
			name: "update lowest priority to highest",
			cq:   NewConnectionQueue(),
			args: args{
				lastUpdate: time.Date(1970, 0, 0, 0, 0, 0, 0, time.Now().Location()),
			},
			want: &models.Connection{
				Conn:       nil,
				LastUpdate: time.Date(1970, 0, 0, 0, 0, 0, 0, time.Now().Location()),
				Index:      -1,
			},
		},
	}
	tests[0].cq.Push(&models.Connection{
		Conn:       nil,
		LastUpdate: time.Date(1990, 0, 0, 0, 0, 0, 0, time.Now().Location()),
		Index:      0,
	})
	tests[0].cq.Push(&models.Connection{
		Conn:       nil,
		LastUpdate: time.Date(2000, 0, 0, 0, 0, 0, 0, time.Now().Location()),
		Index:      0,
	})
	tests[0].cq.Push(&models.Connection{
		Conn:       nil,
		LastUpdate: time.Date(2020, 0, 0, 0, 0, 0, 0, time.Now().Location()),
		Index:      0,
	})

	for _, tt := range tests {
		lowest := tt.cq[tt.cq.Len()-1]
		tt.cq.Update(lowest, tt.args.lastUpdate)
		got := tt.cq.Pop()
		assert.Equal(t, tt.want, got)
	}
}

func TestNewConnectionQueue(t *testing.T) {
	tests := []struct {
		name string
		want ConnectionQueue
	}{
		{
			name: "empty queue",
			want: ConnectionQueue{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewConnectionQueue()
			assert.Equal(t, tt.want, got)
		})
	}
}
