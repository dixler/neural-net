package neural

import (
   "fmt"
   "math/rand"
)

type NeuralNet struct {
   root *layer
   activation func([]int) int
}

type unit struct {
   weights []int
   bias int
}

type layer struct {
   units []unit
   next *layer
}

func random() int {
   num := rand.Int()
   isNegative := (rand.Int() % 2) == 0
   if isNegative {
      num *= -1
   }
   return num
}

func GenerateNN(blueprint []int, activation func([]int) int) NeuralNet {
   // this code generates a random neural network
   // with blueprint according to a
   root := NeuralNet{root: nil}
   root.activation = activation
   num_inputs := blueprint[0]
   var last_layer *layer
   var cur_layer *layer

   for _, num_units := range blueprint {
      cur_layer = &layer{}
      cur_layer.units = make([]unit, num_units, num_units)
      units := cur_layer.units

      // build layer
      for i := 0; i < num_units; i++ {

         // build node
         var cur_node *unit = &units[i]
         cur_node.weights = make([]int, num_inputs, num_inputs)
         cur_node.bias = random()
         // set weights
         for i := 0; i < num_inputs; i++ {
            cur_node.weights[i] = random()
         }
      }
      num_inputs = num_units
      if last_layer == nil{
         root.root = cur_layer
         last_layer = cur_layer
         continue
      }
      last_layer.next = cur_layer
      last_layer = cur_layer
   }

   cur_layer.next = nil
   return root
}

func (n *NeuralNet) DumpNN() {
   for layer := n.root; layer != nil; layer = layer.next {
      for _, node := range layer.units {
         fmt.Print(node.weights[0], " ")
      }
      fmt.Println("")
   }
}

func (n *NeuralNet) Process(input []int) int {
   for layer := n.root; layer != nil; layer = layer.next {
      output := make([]int , len(layer.units), len(layer.units)) // optimizeable

      for i, node := range layer.units {
         output[i] = dotProduct(input, node.weights) + node.bias
      }

      input = output
   }
   return n.activation(input)
}

func dotProduct(a, b []int) int {
   dot := 0
   for i, _ := range(a) {
      dot += (a[i]*b[i])
   }
   return dot
}

const PRECISION = 10000
func inherit(feature_a, feature_b int) int {
   if rand.Int() % 3 > 0 {
      var a_likeness float64 = rand.Float64()/2.0+0.5
      var b_likeness float64 = 1.0-a_likeness
      return int(float64(feature_a)*a_likeness+float64(feature_b)*b_likeness)
   } else {
      if rand.Int() % 2 == 0 {
         return feature_a
      }
      return feature_b
   }
}

func breed(a, b unit) unit {
   child := unit{}
   child.weights = make([]int, len(a.weights), len(a.weights))
   child.bias = inherit(a.bias, b.bias)
   for i, _ := range a.weights {
      child.weights[i] = inherit(a.weights[i], b.weights[i])
   }
   return child
}

func Breed(a, b NeuralNet) NeuralNet {
   layer_a, layer_b := a.root, b.root

   child := NeuralNet {}
   var prev_layer_child *layer = nil

   for ; layer_a != nil; layer_a, layer_b=layer_a.next, layer_b.next {
      layer_child := &layer{}
      layer_child.units = make([]unit, len(layer_a.units), len(layer_a.units))

      // build child layer
      for i, _ := range layer_a.units {
         layer_child.units[i] = breed(layer_a.units[i], layer_b.units[i])
      }

      if prev_layer_child == nil{
         child.root = layer_child
         prev_layer_child = layer_child
         continue
      }

      prev_layer_child.next = layer_child
      prev_layer_child = layer_child
   }

   // handle activator

   child.activation = b.activation
   if rand.Int() % 2 == 0 {
      child.activation = a.activation
   }

   return child
}
