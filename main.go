package main

import (
   "fmt"
   "math/rand"
   "time"
   "./neural"
   "./pong"
)
func genBlueprint(num_layers, inputs, outputs int) []int {
   blueprint := make([]int, num_layers, num_layers)
   blueprint[0] = inputs
   for i := 1; i < num_layers-1; i++ {
      blueprint[i] = 1
   }
   blueprint[len(blueprint)-1] = outputs
   return blueprint
}

func testWorld() {
   type pair struct {
      score int
      bounces int
      nn neural.NeuralNet
   }
   num_networks := 2000
   networks := make([]neural.NeuralNet, num_networks, num_networks)
   blueprint := []int{10, 10}
   activation := func(input []int) int {
      max := input[0]
      idx := 0
      for i, cur := range input {
         if cur > max {
            max = cur
            idx = i
         }
      }
      return idx
   }
   ch := make(chan pair)

   // monitoring:
   first := true
   last_score := 0
   for i, _ := range networks {
      go func(i int) {
         networks[i] = neural.GenerateNN(blueprint, activation)
         ch<-pair{}
      }(i)
   }
   for i := 0; i < len(networks); i++ {
      <-ch
   }

   for i:=0;; i++{
      for _, network := range networks {
         go func(nn neural.NeuralNet) {
            w := pong.GenWorld(100, 200)
            for ;; {
               input := w.GetState()
               filtered_input := []int{input[1], input[5]}
               decision := nn.Process(filtered_input)
               //fmt.Println("decision:", decision)
               if !w.Tick(decision) {
                  ch<-pair{score: w.Score, bounces: w.Bounces, nn: nn}
                  return
               }
            }
         }(network)
      }
      cur_result := <-ch
      best_score, best_nn, best_bounces := cur_result.score, cur_result.nn, cur_result.bounces
      for i := 1; i < num_networks; i++ {
         cur_result := <-ch
         cur_score, cur_nn, cur_bounces :=  cur_result.score, cur_result.nn, cur_result.bounces
         if cur_score > best_score {
            best_score, best_nn, best_bounces = cur_score, cur_nn, cur_bounces
         }
      }

      // metrics
      if first {
         last_score = best_score
         fmt.Println(i, ":", best_bounces, best_score)
         first = false
      } else if best_score > last_score {
         last_score = best_score
         fmt.Println(i, ":", best_bounces, best_score)
      }
      // breeding function
      networks[0] = best_nn
      for i := num_networks/10; i < num_networks; i++{
         go func(i int) {
            networks[i] = neural.GenerateNN(blueprint, activation)
            ch<-pair{}
         }(i)
      }
      for i := num_networks/10; i < num_networks; i++{
         <-ch
      }



      for i := 1; i < num_networks; i++{
         go func(i int) {
            networks[i] = neural.Breed(best_nn, networks[i])
            ch<-pair{}
         }(i)
      }
      for i := 1; i < num_networks; i++{
         <-ch
      }
   }

   return
}

func main() {
   rand.Seed(time.Now().UTC().UnixNano())
   //testBreed()
   //testDecision()
   testWorld()
   return
}
