package main

import (
	"fmt"
	"image/color"
	"log"
	"math"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const W = 720.0
const H = 480.0
const SPATIAL_W = 64
const SPATIAL_H = 36
const REPEL_R = W / SPATIAL_W
const NODE_COUNT = 4
const SPRING_COUNT = 3
const DT = 0.001
const DRAG = 0.99

type Game struct {
	nodes       [NODE_COUNT]node
	spatial_map [SPATIAL_W][SPATIAL_H][]vec2
	springs     [SPRING_COUNT]spring
}
type vec2 struct {
	X float32
	Y float32
}

type node struct {
	pos vec2
	vel vec2
}

type node_id int

type spring struct {
	l node_id
	r node_id
}

func use(x any) {
	if false {
		use(x)
	}
}

func (v vec2) Add(other vec2) vec2 {
	return vec2{
		X: v.X + other.X,
		Y: v.Y + other.Y,
	}
}

func (v vec2) Sub(other vec2) vec2 {
	return vec2{
		X: v.X - other.X,
		Y: v.Y - other.Y,
	}
}

func (v vec2) Mul(scalar float32) vec2 {
	return vec2{
		X: v.X * scalar,
		Y: v.Y * scalar,
	}
}

func (v vec2) Div(scalar float32) vec2 {
	return v.Mul(1 / scalar)
}

func (v vec2) Mag() float32 {
	return float32(math.Sqrt(float64(v.X*v.X + v.Y*v.Y)))
}

func (v vec2) Norm() vec2 {
	mag := v.Mag()
	if mag > 0.00000001 {
		return v.Div(v.Mag())
	}
	return vec2{X: 0, Y: 0}
}

func (v vec2) WithMag(mag float32) vec2 {
	return v.Norm().Mul(mag)
}

func (p vec2) String() string {
	return fmt.Sprintf("(%d, %d)", p.X, p.Y)
}

func clamp(v, lower, upper int) int {
	return min(upper, max(lower, v))
}

func (g *Game) Update() error {
	for i := range g.springs {
		spring := g.springs[i]
		n1 := &g.nodes[spring.l]
		n2 := &g.nodes[spring.r]
		d := n1.pos.Sub(n2.pos)
		n1.vel = n1.vel.Sub(d)
		n2.vel = n2.vel.Add(d)
	}
	for x := range g.spatial_map {
		for y := range g.spatial_map[x] {
			g.spatial_map[x][y] = make([]vec2, 0)
		}
	}
	for i := range g.nodes {
		node := &g.nodes[i]
		pos := node.pos
		sx := clamp(int(pos.X/W), 0, SPATIAL_W-1)
		sy := clamp(int(pos.Y/H), 0, SPATIAL_H-1)
		g.spatial_map[sx][sy] = append(g.spatial_map[sx][sy], node.pos)
	}
	for i := range g.nodes {
		node := &g.nodes[i]
		pos := node.pos
		sx := clamp(int(pos.X/W), 0, SPATIAL_W-1)
		sy := clamp(int(pos.Y/H), 0, SPATIAL_H-1)
		for _dx := range 3 {
			for _dy := range 3 {
				dx := _dx - 1
				dy := _dy - 1
				px := sx + dx
				py := sy + dy
				if px < 0 || px >= SPATIAL_W || py < 0 || py >= SPATIAL_H {
					continue
				}
				spatial_group := g.spatial_map[px][py]
				for i := range spatial_group {
					gn := spatial_group[i]
					diff := node.pos.Sub(gn)
					//use(diff)
					fmt.Println(diff, diff.Mag())
					node.vel = node.vel.Add(diff.WithMag(REPEL_R - diff.Mag()).Mul(-0.1))
				}
			}
		}
		node.vel = node.vel.Add(vec2{X: W / 2, Y: H / 2}.Sub(node.pos).Mul(2))
		node.pos = node.pos.Add(node.vel.Mul(DT))
		node.vel = node.vel.Mul(DRAG)
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	for i := range g.springs {
		spring := g.springs[i]
		p1 := g.nodes[spring.l].pos
		p2 := g.nodes[spring.r].pos
		vector.StrokeLine(screen, p1.X, p1.Y, p2.X, p2.Y, 1, color.White, false)
	}
}

func (g *Game) Layout(outsideWidth int, outsideHeight int) (screenWidth int, screenHeight int) {
	return W, H
}

func main() {
	nodes := new([NODE_COUNT]node)
	springs := new([SPRING_COUNT]spring)
	rng := rand.New(rand.NewSource(time.Now().Unix()))
	for i := range NODE_COUNT {
		nodes[i].pos = vec2{X: W * rng.Float32(), Y: H * rng.Float32()}
	}
	for i := range SPRING_COUNT {
		/*springs[i].l = node_id(rng.Int31() % NODE_COUNT)
		springs[i].r = node_id(rng.Int31() % NODE_COUNT)*/
		springs[i].l = node_id(0)
		springs[i].r = node_id(i + 1)
	}
	ebiten.SetWindowSize(W, H)
	ebiten.SetWindowTitle("Hello, World!")
	if err := ebiten.RunGame(&Game{nodes: *nodes, springs: *springs}); err != nil {
		log.Fatal(err)
	}
}