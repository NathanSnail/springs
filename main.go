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
const RENDER_SCALE = 2.0
const SPATIAL_W = 32
const SPATIAL_H = 18
const SPATIAL_RAD = 2
const SPATIAL_SIZE = W / SPATIAL_W
const REPEL_R = SPATIAL_SIZE * SPATIAL_RAD
const REPEL_FORCE = 3
const CENTRE_FORCE = 0.03
const CURSOR_FORCE = -10000
const NODE_COUNT = 64
const SPRING_COUNT = 64
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
	lmb := ebiten.IsMouseButtonPressed(ebiten.MouseButton0)
	cx, cy := ebiten.CursorPosition()
	cpos := vec2{X: float32(cx) / RENDER_SCALE, Y: float32(cy) / RENDER_SCALE}
	fmt.Println(cpos)
	for i := range g.nodes {
		node := &g.nodes[i]
		pos := node.pos
		sx := clamp(int(pos.X/W), 0, SPATIAL_W-1)
		sy := clamp(int(pos.Y/H), 0, SPATIAL_H-1)
		for _dx := range SPATIAL_RAD*2 + 1 {
			for _dy := range SPATIAL_RAD*2 + 1 {
				dx := _dx - SPATIAL_RAD
				dy := _dy - SPATIAL_RAD
				px := sx + dx
				py := sy + dy
				if px < 0 || px >= SPATIAL_W || py < 0 || py >= SPATIAL_H {
					continue
				}
				spatial_group := g.spatial_map[px][py]
				for i := range spatial_group {
					gn := spatial_group[i]
					diff := node.pos.Sub(gn)
					force := diff.WithMag(max(0, (REPEL_R-diff.Mag())*REPEL_FORCE))
					node.vel = node.vel.Add(force)
				}
			}
		}
		node.vel = node.vel.Add(vec2{X: W / 2, Y: H / 2}.Sub(node.pos).Mul(CENTRE_FORCE))
		if lmb {
			delta := cpos.Sub(node.pos)
			mag := min(1, 1/delta.Mag())
			delta = delta.WithMag(mag * mag)
			node.vel = node.vel.Add(delta.Mul(CURSOR_FORCE))
		}
		node.pos = node.pos.Add(node.vel.Mul(DT))
		node.vel = node.vel.Mul(DRAG)
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	for i := range g.springs {
		spring := g.springs[i]
		p1 := g.nodes[spring.l].pos.Mul(RENDER_SCALE)
		p2 := g.nodes[spring.r].pos.Mul(RENDER_SCALE)
		vector.StrokeLine(screen, p1.X, p1.Y, p2.X, p2.Y, 1, color.White, false)
	}
}

func (g *Game) Layout(outsideWidth int, outsideHeight int) (screenWidth int, screenHeight int) {
	return W * RENDER_SCALE, H * RENDER_SCALE
}

func main() {
	nodes := new([NODE_COUNT]node)
	springs := new([SPRING_COUNT]spring)
	rng := rand.New(rand.NewSource(time.Now().Unix()))
	for i := range NODE_COUNT {
		nodes[i].pos = vec2{X: W * rng.Float32(), Y: H * rng.Float32()}
	}
	for i := range SPRING_COUNT {
		// /*
		springs[i].l = node_id(rng.Int31() % NODE_COUNT)
		springs[i].r = node_id(rng.Int31() % NODE_COUNT)
		// */
		/*
			springs[i].l = node_id(0)
			springs[i].r = node_id(i + 1)
		*/
	}
	ebiten.SetWindowSize(W*RENDER_SCALE, H*RENDER_SCALE)
	ebiten.SetWindowTitle("Spring Toy")
	if err := ebiten.RunGame(&Game{nodes: *nodes, springs: *springs}); err != nil {
		log.Fatal(err)
	}
}
