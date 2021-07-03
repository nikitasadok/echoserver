package connectionQueue

import (
	"echoServer/models"
	"github.com/stretchr/testify/assert"
	"reflect"
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
			if got := tt.cq.Len(); got != tt.want {
				t.Errorf("Len() = %v, want %v", got, tt.want)
			}
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
			if got := tt.cq.Less(tt.args.i, tt.args.j); got != tt.want {
				t.Errorf("Less() = %v, want %v", got, tt.want)
			}
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
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
		})
	}
}

func TestConnectionQueue_Swap(t *testing.T) {
	type args struct {
		i int
		j int
	}
	tests := []struct {
		name string
		cq   ConnectionQueue
		args args
	}{
		{
			name: "simple swap",
			cq:   NewConnectionQueue(),
			args: args{
				i: 0,
				j: 1,
			},
		},
	}

	tests[0].cq.Push(&models.Connection{
		Conn:       nil,
		LastUpdate: time.Time{},
		Index:      0,
	})

	// tests[0].cq.Push()
	for _, tt := range tests {
		assert.Equal(t, tt.cq[tt.args.i], tt.cq[tt.args.j])
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
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
		})
	}
}

func TestNewConnectionQueue(t *testing.T) {
	tests := []struct {
		name string
		want ConnectionQueue
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewConnectionQueue(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewConnectionQueue() = %v, want %v", got, tt.want)
			}
		})
	}
}
