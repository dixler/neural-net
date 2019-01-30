package main

import (
   "fmt"
   "math/rand"
   "time"
   "./neural"
   "./pong"
)
func testDecision() {
   blueprint := []int{6, 4, 4, 3}

   for _, val := range blueprint {
      fmt.Println(val)
   }

   nn := neural.GenerateNN(blueprint,
      func(input []int) int {
         max := 0
         idx := 0
         for i, val := range input {
            if val > max {
               max = val
               idx = i
            }
         }
         return idx
      })
   nn.DumpNN()
   fmt.Println("result", nn.Process([]int{3, 3, 3, 3, 3, 3}))
   return
}
func genBlueprint(num_layers, inputs, outputs int) []int {
   blueprint := make([]int, num_layers, num_layers)
   blueprint[0] = inputs
   for i := 1; i < num_layers-1; i++ {
      blueprint[i] = 6
   }
   blueprint[len(blueprint)-1] = outputs
   return blueprint
}
func testBreed() {
   // Randomly generated NN
   num_layers := 100
   blueprint := genBlueprint(num_layers, 6, 3)
   activation := func(input []int) int {
         max := 0
         idx := 0
         for i, val := range input {
            if val > max {
               max = val
               idx = i
            }
         }
         return idx
      }
   mismatches := 0
   nn1 := neural.GenerateNN(blueprint, activation)
   nn2 := neural.GenerateNN(blueprint, activation)

   res_1 := nn1.Process([]int{1,1,1,1,1,1})

   child := neural.Breed(nn1, nn2)

   for i := 0; i < 10000; i++ {
      res_child := child.Process([]int{1,1,1,1,1,1})


      if res_1 != res_child {
         fmt.Println("mismatch", i, res_1, res_child)
         child = neural.Breed(nn1, child)
         mismatches++
      }
      child = neural.Breed(nn1, child)
   }
   fmt.Println("mismatches:", mismatches)



}

func testWorld() {
   type pair struct {
      score int
      nn neural.NeuralNet
   }
   num_networks := 1000
   networks := make([]neural.NeuralNet, num_networks, num_networks)
   blueprint := genBlueprint(3, 6, 3)
   activation := func(input []int) int {
      max := 0
      idx := 0
      for i, val := range input {
         if val > max {
            max = val
            idx = i
         }
      }
      return idx
   }
   ch := make(chan pair)

   for i, _ := range networks {
      networks[i] = neural.GenerateNN(blueprint, activation)
   }

   for i:=0;; i++{
      for _, network := range networks {
         go func(nn neural.NeuralNet) {
            w := pong.GenWorld(100, 200)
            for ;; {
               input := w.GetState()
               decision := nn.Process(input)
               //fmt.Println("decision:", decision)
               if !w.Tick(decision) {
                  ch<-pair{score: w.Score, nn: nn}
                  return
               }
            }
         }(network)
      }
      best_score, best_nn := -1, neural.NeuralNet{}
      for i := 0; i < num_networks; i++ {
         cur_result := <-ch
         cur_score, cur_nn :=  cur_result.score, cur_result.nn
         if cur_score > best_score {
            best_score, best_nn = cur_score, cur_nn
         }
      }

      fmt.Println(i, ":", best_score)
      // breeding function
      networks[0] = best_nn
      for i := 1; i < num_networks/10; i++{
         networks[i] = neural.Breed(networks[i], best_nn)
      }
      for i := num_networks/10; i < num_networks; i++{
         networks[i] = neural.GenerateNN(blueprint, activation)
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
