package pong
import (
   //"math"
)

type object struct {
   x        int
   y        int
   dx       int
   dy       int
   width    int
   height   int
}

type World struct {
   paddle   *object
   ball     *object
   width    int
   height   int
   Score    int
   Bounces  int
}

const (
   Up    int = iota
   NoOp
   Down
)

func GenWorld(height, width int) World {
   world := World{ width: width, height:  height, Score: 0}
   world.paddle = &object{ x: 0, y: (height-(height/5))/2, dx: 0, dy: 0, width: 1, height: height/5, }
   world.ball = &object{x: width/2, y: height/2, dx: 1, dy: 1, width: 1, height: 1}
   return world
}

func collision(pad, ball *object) bool {
   // handle collision on x bounds
   if ball.x == pad.x+1 {
      // handle collision on y bounds
      if ball.y >= pad.y && ball.y < pad.y+pad.height {
         // in between pad edges
         return true
      }
   }
   return false
}

func (w *World) Tick (action int) bool {
   deltaY := 0
   switch action {
      case Up:
         deltaY-=1
      case Down:
         deltaY+=1
   }

   if deltaY < 0 || deltaY+w.paddle.height > w.height {
      deltaY = 0
   }

   w.paddle.y += deltaY
   if collision(w.paddle, w.ball) {
      w.ball.dx *= -1
      w.Bounces++
      w.Score += 10000000
   }
   if w.ball.x  > w.width{
      w.ball.dx *= -1
   }

   if w.ball.y < 0 || w.ball.y > w.height{
      w.ball.dy *= -1
   }
   if w.ball.x < 0 {
      return false
   }

   _, pad_midpoint_y := w.paddle.x, w.paddle.y+w.paddle.height/2

   _, w.paddle.y = w.paddle.x+w.paddle.dx, w.paddle.y+w.paddle.dy
   w.ball.x, w.ball.y = w.ball.x+w.ball.dx, w.ball.y+w.ball.dy

   //distance := math.Pow(float64(w.ball.x - pad_midpoint_x), 2.0) + math.Pow(float64(w.ball.y - pad_midpoint_y), 2.0)
   distance := -((w.ball.y - pad_midpoint_y)*(w.ball.y - pad_midpoint_y))

   w.Score += distance
   return true
}

func (w World) GetState () []int {
   return []int{
      w.ball.x,
      w.ball.y,
      w.ball.dx,
      w.ball.dy,
      w.paddle.x,
      w.paddle.y }
}
