package pbm

import (
	"fmt"
	"time"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ErrorCursor struct {
	cerr error
}

func (c ErrorCursor) Error() string {
	return fmt.Sprintln("cursor was closed with:", c.cerr)
}

func (p *PBM) ListenCmd() (<-chan Cmd, <-chan error, error) {
	cmd := make(chan Cmd)
	errc := make(chan error)
	go func() {
		defer close(cmd)
		defer close(errc)
		// defer cur.Close(p.ctx)
		ts := time.Now().UTC().Unix()
		for {
			cur, err := p.Conn.Database(DB).Collection(CmdStreamCollection).Find(
				p.ctx,
				bson.M{"ts": bson.M{"$gte": ts}},
				options.Find(),
			)
			if err != nil {
				errc <- errors.Wrap(err, "watch the cmd stream")
				return
			}

			for cur.Next(p.ctx) {
				c := Cmd{}
				err := cur.Decode(&c)
				if err != nil {
					errc <- errors.Wrap(err, "message decode")
					continue
				}

				cmd <- c
				ts = time.Now().UTC().Unix()
			}
			if cur.Err() != nil {
				errc <- ErrorCursor{cerr: cur.Err()}
				cur.Close(p.ctx)
				return
			}
			cur.Close(p.ctx)
			time.Sleep(time.Second * 1)
		}
	}()

	return cmd, errc, nil
}

func (p *PBM) SendCmd(cmd Cmd) error {
	cmd.TS = time.Now().UTC().Unix()
	_, err := p.Conn.Database(DB).Collection(CmdStreamCollection).InsertOne(p.ctx, cmd)
	return err
}
