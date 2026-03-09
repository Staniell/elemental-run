package game

import (
	"fmt"
	"image/color"
	"math"
	"math/rand"
	"time"

	"github.com/hajimehoshi/bitmapfont/v4"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	ebitext "github.com/hajimehoshi/ebiten/v2/text"
)

const (
	ScreenWidth  = 480
	ScreenHeight = 270

	groundY       = 220
	playerScreenX = 112

	gravity      = 0.44
	jumpVelocity = -8.2
	maxFallSpeed = 9.2

	baseRunSpeed   = 3.1
	coyoteFrames   = 8
	jumpBufferSpan = 8

	maxHealth = 5

	elementCount = 3
	introFrames  = 540
)

type element int

const (
	fire element = iota
	ice
	thunder
)

type enemyKind int

const (
	enemyGround enemyKind = iota
	enemyTurret
	enemyFlyer
)

type tileKind int

const (
	tileGround tileKind = iota
	tilePlatform
	tileCrate
	tileSpike
)

type gameState int

const (
	stateMenu gameState = iota
	statePlaying
	stateGameOver
)

const (
	menuStart = iota
	menuQuit
)

var menuLabels = [...]string{"Start", "Quit"}

var elementNames = [...]string{"Fire", "Ice", "Thunder"}

var hudSmallFont = bitmapfont.Face

var elementColors = [...]color.RGBA{
	{R: 255, G: 141, B: 67, A: 255},
	{R: 126, G: 230, B: 255, A: 255},
	{R: 247, G: 233, B: 92, A: 255},
}

type rect struct {
	X float64
	Y float64
	W float64
	H float64
}

func (r rect) intersects(o rect) bool {
	return r.X < o.X+o.W && r.X+r.W > o.X && r.Y < o.Y+o.H && r.Y+r.H > o.Y
}

func (r rect) overlapsX(o rect) bool {
	return r.X < o.X+o.W && r.X+r.W > o.X
}

type player struct {
	y           float64
	vy          float64
	w           float64
	h           float64
	onGround    bool
	coyote      int
	jumpBuffer  int
	invuln      int
	hp          int
	attackFlash int
	lastCast    element
	cooldowns   [elementCount]int
}

func (p player) rect() rect {
	return rect{X: playerScreenX, Y: p.y, W: p.w, H: p.h}
}

type platform struct {
	x    float64
	y    float64
	w    float64
	h    float64
	tile tileKind
}

func (p platform) rect() rect {
	return rect{X: p.x, Y: p.y, W: p.w, H: p.h}
}

type hazard struct {
	x    float64
	y    float64
	w    float64
	h    float64
	tile tileKind
}

func (h hazard) rect() rect {
	return rect{X: h.x, Y: h.y, W: h.w, H: h.h}
}

type enemy struct {
	id         int
	kind       enemyKind
	element    element
	x          float64
	y          float64
	baseY      float64
	w          float64
	h          float64
	hp         int
	maxHP      int
	shootTimer int
	phase      float64
	flash      int
	stun       int
	freeze     int
	burn       int
	dead       bool
}

func (e enemy) rect() rect {
	return rect{X: e.x, Y: e.y, W: e.w, H: e.h}
}

type projectile struct {
	x             float64
	y             float64
	vx            float64
	vy            float64
	w             float64
	h             float64
	damage        int
	ttl           int
	hitsRemaining int
	element       element
	fromPlayer    bool
	dead          bool
}

func (p projectile) rect() rect {
	return rect{X: p.x, Y: p.y, W: p.w, H: p.h}
}

type swing struct {
	element element
	ttl     int
	damage  int
	w       float64
	h       float64
	xOff    float64
	yOff    float64
	hit     map[int]bool
}

func (s swing) rect(p player) rect {
	return rect{X: playerScreenX + s.xOff, Y: p.y + s.yOff, W: s.w, H: s.h}
}

type Game struct {
	assets *gameAssets
	rng    *rand.Rand

	player      player
	platforms   []platform
	hazards     []hazard
	enemies     []enemy
	projectiles []projectile
	swings      []swing

	nextEnemyID int
	nextChunkX  float64

	runSpeed  float64
	distance  float64
	score     int
	kills     int
	ticks     int
	bestScore int

	shake       int
	state       gameState
	menuIndex   int
	introFrames int
	hudFocus    int
	hudAbility  float64
	savePath    string

	bgFarOffset  float64
	bgMidOffset  float64
	bgNearOffset float64
}

func New() (*Game, error) {
	assets, err := loadAssets()
	if err != nil {
		return nil, err
	}

	g := &Game{
		assets:   assets,
		rng:      rand.New(rand.NewSource(time.Now().UnixNano())),
		savePath: defaultSavePath(),
	}
	g.loadProgress()
	g.openMenu()
	return g, nil
}

func (g *Game) startRun() {
	g.player = player{
		y:        groundY - 40,
		w:        30,
		h:        40,
		onGround: true,
		hp:       maxHealth,
		lastCast: fire,
	}

	g.platforms = nil
	g.hazards = nil
	g.enemies = nil
	g.projectiles = nil
	g.swings = nil
	g.nextEnemyID = 1
	g.nextChunkX = ScreenWidth + 64
	g.runSpeed = baseRunSpeed
	g.distance = 0
	g.score = 0
	g.kills = 0
	g.ticks = 0
	g.shake = 0
	g.state = statePlaying
	g.introFrames = introFrames
	g.hudFocus = 120
	g.hudAbility = 0
	g.bgFarOffset = 0
	g.bgMidOffset = 0
	g.bgNearOffset = 0

	for g.nextChunkX < ScreenWidth+420 {
		g.spawnChunk()
	}
}

func (g *Game) openMenu() {
	g.setBestScore(g.score)
	g.startRun()
	g.state = stateMenu
	g.menuIndex = menuStart
	g.introFrames = 0
	g.hudFocus = 0
	g.hudAbility = 0
	g.shake = 0
}

func (g *Game) Layout(_, _ int) (int, int) {
	return ScreenWidth, ScreenHeight
}

func (g *Game) Update() error {
	g.ticks++
	if g.shake > 0 {
		g.shake--
	}

	switch g.state {
	case stateMenu:
		return g.updateMenu()
	case statePlaying:
		return g.updatePlaying()
	case stateGameOver:
		return g.updateGameOver()
	default:
		return nil
	}
}

func (g *Game) updateMenu() error {
	g.bgFarOffset = wrapLayer(g.bgFarOffset - 0.18)
	g.bgMidOffset = wrapLayer(g.bgMidOffset - 0.34)
	g.bgNearOffset = wrapLayer(g.bgNearOffset - 0.52)

	if inpututil.IsKeyJustPressed(ebiten.KeyArrowUp) || inpututil.IsKeyJustPressed(ebiten.KeyW) {
		g.moveMenuSelection(-1)
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyArrowDown) || inpututil.IsKeyJustPressed(ebiten.KeyS) {
		g.moveMenuSelection(1)
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		return ebiten.Termination
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) || inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		switch g.menuIndex {
		case menuStart:
			g.startRun()
		case menuQuit:
			return ebiten.Termination
		}
	}

	return nil
}

func (g *Game) updatePlaying() error {
	if g.introFrames > 0 {
		g.introFrames--
	}
	if g.hudFocus > 0 {
		g.hudFocus--
	}
	g.updateAbilityHUDAnimation()
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		g.openMenu()
		return nil
	}

	g.handlePlayInput()
	g.advanceDifficulty()
	g.scrollBackgrounds()
	g.scrollTerrain()
	g.spawnChunksIfNeeded()
	g.updatePlayer()
	g.updateSwings()
	g.updateEnemies()
	g.updateProjectiles()
	g.resolveCombat()
	g.resolveThreats()
	g.cleanup()
	g.score = int(g.distance/10) + g.kills*35

	return nil
}

func (g *Game) updateGameOver() error {
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		g.openMenu()
		return nil
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyR) || inpututil.IsKeyJustPressed(ebiten.KeyEnter) || inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		g.startRun()
	}
	return nil
}

func (g *Game) handlePlayInput() {
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) || inpututil.IsKeyJustPressed(ebiten.KeyArrowUp) || inpututil.IsKeyJustPressed(ebiten.KeyW) {
		g.player.jumpBuffer = jumpBufferSpan
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyJ) || inpututil.IsKeyJustPressed(ebiten.KeyDigit1) {
		g.tryCast(fire)
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyK) || inpututil.IsKeyJustPressed(ebiten.KeyDigit2) {
		g.tryCast(ice)
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyL) || inpututil.IsKeyJustPressed(ebiten.KeyDigit3) {
		g.tryCast(thunder)
	}
}

func (g *Game) moveMenuSelection(delta int) {
	g.menuIndex = (g.menuIndex + delta + len(menuLabels)) % len(menuLabels)
}

func (g *Game) advanceDifficulty() {
	g.distance += g.runSpeed
	g.runSpeed = baseRunSpeed + math.Min(2.0, g.distance/2600)
}

func (g *Game) scrollBackgrounds() {
	g.bgFarOffset = wrapLayer(g.bgFarOffset - g.runSpeed*0.15)
	g.bgMidOffset = wrapLayer(g.bgMidOffset - g.runSpeed*0.35)
	g.bgNearOffset = wrapLayer(g.bgNearOffset - g.runSpeed*0.65)
}

func (g *Game) scrollTerrain() {
	g.nextChunkX -= g.runSpeed
	for i := range g.platforms {
		g.platforms[i].x -= g.runSpeed
	}
	for i := range g.hazards {
		g.hazards[i].x -= g.runSpeed
	}
}

func (g *Game) spawnChunksIfNeeded() {
	for g.nextChunkX < ScreenWidth+320 {
		g.spawnChunk()
	}
}

func (g *Game) updatePlayer() {
	p := &g.player

	if p.invuln > 0 {
		p.invuln--
	}
	if p.attackFlash > 0 {
		p.attackFlash--
	}
	for i := 0; i < elementCount; i++ {
		if p.cooldowns[i] > 0 {
			p.cooldowns[i]--
		}
	}
	if p.jumpBuffer > 0 {
		p.jumpBuffer--
	}
	if p.onGround {
		p.coyote = coyoteFrames
	} else if p.coyote > 0 {
		p.coyote--
	}

	if p.jumpBuffer > 0 && (p.onGround || p.coyote > 0) {
		p.vy = jumpVelocity
		p.onGround = false
		p.coyote = 0
		p.jumpBuffer = 0
	}

	prevY := p.y
	prevTop := prevY
	prevBottom := prevY + p.h

	p.vy += gravity
	if p.vy > maxFallSpeed {
		p.vy = maxFallSpeed
	}
	p.y += p.vy

	landY := groundY - p.h
	landed := false
	if p.vy >= 0 && prevBottom <= groundY && p.y+p.h >= groundY {
		landed = true
	}

	pRect := p.rect()
	for _, plat := range g.platforms {
		pr := plat.rect()
		if !pRect.overlapsX(pr) {
			continue
		}

		if p.vy >= 0 && prevBottom <= plat.y && p.y+p.h >= plat.y {
			candidate := plat.y - p.h
			if !landed || candidate < landY {
				landY = candidate
				landed = true
			}
			continue
		}

		if p.vy < 0 && prevTop >= plat.y+plat.h && p.y <= plat.y+plat.h {
			p.y = plat.y + plat.h
			p.vy = 0
		}
	}

	if landed {
		p.y = landY
		p.vy = 0
		p.onGround = true
	} else {
		p.onGround = false
	}
}

func (g *Game) updateSwings() {
	for i := range g.swings {
		if g.swings[i].ttl > 0 {
			g.swings[i].ttl--
		}
	}
}

func (g *Game) updateEnemies() {
	for i := range g.enemies {
		e := &g.enemies[i]
		if e.dead {
			continue
		}

		if e.flash > 0 {
			e.flash--
		}
		if e.freeze > 0 {
			e.freeze--
		}
		if e.stun > 0 {
			e.stun--
		}
		if e.burn > 0 {
			e.burn--
			if e.burn%18 == 0 {
				e.hp--
				e.flash = 4
				if e.hp <= 0 {
					e.dead = true
					g.kills++
					g.shake = maxInt(g.shake, 4)
					continue
				}
			}
		}

		speedScale := 1.0
		if e.freeze > 0 {
			speedScale = 0.5
		}
		if e.stun > 0 {
			speedScale = 0.0
		}

		switch e.kind {
		case enemyGround:
			e.x -= g.runSpeed * speedScale
		case enemyTurret:
			e.x -= g.runSpeed * speedScale
			if e.stun == 0 {
				e.shootTimer--
				if e.shootTimer <= 0 && e.x > playerScreenX+32 && e.x < ScreenWidth+48 {
					g.spawnEnemyProjectile(e)
					e.shootTimer = 84 + g.rng.Intn(34) - minInt(30, int(g.distance/180))
				}
			}
		case enemyFlyer:
			e.phase += 0.08
			e.x -= g.runSpeed*speedScale + 0.4
			e.y = e.baseY + math.Sin(e.phase)*14
		}
	}
}

func (g *Game) updateProjectiles() {
	for i := range g.projectiles {
		p := &g.projectiles[i]
		if p.dead {
			continue
		}
		p.ttl--
		p.x += p.vx - g.runSpeed
		p.y += p.vy
		if p.ttl <= 0 {
			p.dead = true
		}
	}
}

func (g *Game) resolveCombat() {
	for si := range g.swings {
		s := &g.swings[si]
		if s.ttl <= 0 {
			continue
		}
		sr := s.rect(g.player)

		for pi := range g.projectiles {
			p := &g.projectiles[pi]
			if p.dead || p.fromPlayer {
				continue
			}
			if sr.intersects(p.rect()) {
				p.fromPlayer = true
				p.element = s.element
				p.vx = attackVelocity(s.element) + 1.5
				p.vy = 0
				p.damage = maxInt(p.damage, s.damage+1)
				p.hitsRemaining = maxInt(p.hitsRemaining, 1)
				p.ttl = maxInt(p.ttl, 42)
				g.shake = maxInt(g.shake, 3)
			}
		}

		for ei := range g.enemies {
			e := &g.enemies[ei]
			if e.dead || s.hit[e.id] {
				continue
			}
			if sr.intersects(e.rect()) {
				s.hit[e.id] = true
				g.damageEnemy(e, s.damage, s.element)
			}
		}
	}

	for pi := range g.projectiles {
		p := &g.projectiles[pi]
		if p.dead || !p.fromPlayer {
			continue
		}

		for ei := range g.enemies {
			e := &g.enemies[ei]
			if e.dead {
				continue
			}
			if p.rect().intersects(e.rect()) {
				g.damageEnemy(e, p.damage, p.element)
				p.hitsRemaining--
				if p.hitsRemaining <= 0 {
					p.dead = true
					break
				}
			}
		}
	}
}

func (g *Game) resolveThreats() {
	pRect := g.player.rect()

	for i := range g.platforms {
		plat := g.platforms[i]
		if plat.tile != tileCrate {
			continue
		}
		if pRect.intersects(plat.rect()) && g.player.y+g.player.h > plat.y+6 {
			g.hitPlayer(1)
			break
		}
	}

	if g.state == stateGameOver {
		return
	}

	for i := range g.hazards {
		if pRect.intersects(g.hazards[i].rect()) {
			g.hitPlayer(1)
			break
		}
	}

	if g.state == stateGameOver {
		return
	}

	for i := range g.enemies {
		if g.enemies[i].dead {
			continue
		}
		if pRect.intersects(g.enemies[i].rect()) {
			g.hitPlayer(1)
			break
		}
	}

	if g.state == stateGameOver {
		return
	}

	for i := range g.projectiles {
		p := &g.projectiles[i]
		if p.dead || p.fromPlayer {
			continue
		}
		if pRect.intersects(p.rect()) {
			p.dead = true
			g.hitPlayer(1)
			break
		}
	}
}

func (g *Game) cleanup() {
	platforms := g.platforms[:0]
	for _, p := range g.platforms {
		if p.x+p.w > -64 {
			platforms = append(platforms, p)
		}
	}
	g.platforms = platforms

	hazards := g.hazards[:0]
	for _, h := range g.hazards {
		if h.x+h.w > -64 {
			hazards = append(hazards, h)
		}
	}
	g.hazards = hazards

	enemies := g.enemies[:0]
	for _, e := range g.enemies {
		if !e.dead && e.x+e.w > -80 {
			enemies = append(enemies, e)
		}
	}
	g.enemies = enemies

	projectiles := g.projectiles[:0]
	for _, p := range g.projectiles {
		if p.dead {
			continue
		}
		if p.x+p.w < -32 || p.x > ScreenWidth+48 || p.y > ScreenHeight+48 || p.y+p.h < -32 {
			continue
		}
		projectiles = append(projectiles, p)
	}
	g.projectiles = projectiles

	swings := g.swings[:0]
	for _, s := range g.swings {
		if s.ttl > 0 {
			swings = append(swings, s)
		}
	}
	g.swings = swings
	if g.state == stateGameOver {
		g.setBestScore(g.score)
	}
}

func (g *Game) tryCast(el element) {
	idx := int(el)
	if g.player.cooldowns[idx] > 0 {
		return
	}

	g.player.cooldowns[idx] = attackCooldown(el)
	g.player.lastCast = el
	g.player.attackFlash = 8
	g.hudFocus = maxInt(g.hudFocus, 120)

	w, h, xOff, yOff := attackSwing(el)
	g.swings = append(g.swings, swing{
		element: el,
		ttl:     attackSwingTTL(el),
		damage:  attackDamage(el),
		w:       w,
		h:       h,
		xOff:    xOff,
		yOff:    yOff,
		hit:     map[int]bool{},
	})

	pw, ph := attackProjectileSize(el)
	g.projectiles = append(g.projectiles, projectile{
		x:             playerScreenX + g.player.w - 2,
		y:             g.player.y + attackProjectileSpawnY(el),
		vx:            attackVelocity(el),
		vy:            0,
		w:             pw,
		h:             ph,
		damage:        attackDamage(el),
		ttl:           attackProjectileTTL(el),
		hitsRemaining: attackPierce(el),
		element:       el,
		fromPlayer:    true,
	})
}

func (g *Game) damageEnemy(e *enemy, baseDamage int, attack element) {
	if e.dead {
		return
	}

	damage := baseDamage
	strong := isStrongAgainst(attack, e.element)
	weak := isStrongAgainst(e.element, attack)

	if strong {
		damage = int(math.Round(float64(baseDamage) * 1.8))
	} else if weak {
		damage = maxInt(1, int(math.Round(float64(baseDamage)*0.65)))
	}
	if damage < 1 {
		damage = 1
	}

	e.hp -= damage
	e.flash = 8

	if strong {
		switch attack {
		case fire:
			e.burn = maxInt(e.burn, 90)
		case ice:
			e.freeze = maxInt(e.freeze, 60)
		case thunder:
			e.stun = maxInt(e.stun, 28)
		}
		g.shake = maxInt(g.shake, 6)
	}

	if e.hp <= 0 {
		e.dead = true
		g.kills++
		g.shake = maxInt(g.shake, 5)
	}
}

func (g *Game) hitPlayer(amount int) {
	p := &g.player
	if p.invuln > 0 || g.state == stateGameOver {
		return
	}

	p.hp -= amount
	p.invuln = 50
	g.hudFocus = maxInt(g.hudFocus, 120)
	if p.vy > -4.2 {
		p.vy = -4.2
	}
	g.shake = maxInt(g.shake, 8)

	if p.hp <= 0 {
		p.hp = 0
		g.state = stateGameOver
		g.setBestScore(g.score)
	}
}

func (g *Game) spawnEnemyProjectile(e *enemy) {
	width := 14.0
	height := 14.0
	if e.element == ice {
		width = 18
		height = 14
	}
	g.projectiles = append(g.projectiles, projectile{
		x:             e.x - 4,
		y:             e.y + 8,
		vx:            -2.8,
		vy:            0,
		w:             width,
		h:             height,
		damage:        1,
		ttl:           130,
		hitsRemaining: 1,
		element:       e.element,
		fromPlayer:    false,
	})
}

func (g *Game) spawnChunk() {
	x := g.nextChunkX
	difficulty := int(g.distance / 900)
	roll := g.rng.Intn(100)

	var width float64
	switch {
	case difficulty < 1:
		switch {
		case roll < 35:
			width = g.chunkBreather(x)
		case roll < 68:
			width = g.chunkCrateStep(x)
		default:
			width = g.chunkTurretRoost(x)
		}
	case difficulty < 3:
		switch {
		case roll < 18:
			width = g.chunkBreather(x)
		case roll < 38:
			width = g.chunkCrateStep(x)
		case roll < 58:
			width = g.chunkSpikeLane(x)
		case roll < 80:
			width = g.chunkTurretRoost(x)
		default:
			width = g.chunkSkyBridge(x)
		}
	default:
		switch {
		case roll < 15:
			width = g.chunkCrateStep(x)
		case roll < 35:
			width = g.chunkSpikeLane(x)
		case roll < 58:
			width = g.chunkTurretRoost(x)
		case roll < 80:
			width = g.chunkSkyBridge(x)
		default:
			width = g.chunkMixedGauntlet(x)
		}
	}

	g.nextChunkX += width
}

func (g *Game) chunkBreather(x float64) float64 {
	g.addEnemyOnSurface(enemyGround, x+128, groundY)
	return 180
}

func (g *Game) chunkCrateStep(x float64) float64 {
	g.addPlatform(x+72, groundY-32, 32, 32, tileCrate)
	g.addPlatform(x+128, groundY-64, 32, 32, tileCrate)
	g.addEnemyOnSurface(enemyGround, x+196, groundY)
	return 228
}

func (g *Game) chunkSpikeLane(x float64) float64 {
	g.addHazard(x+82, groundY-18, 32, 18, tileSpike)
	g.addPlatform(x+138, groundY-32, 32, 32, tileCrate)
	g.addEnemyOnSurface(enemyGround, x+204, groundY)
	return 236
}

func (g *Game) chunkTurretRoost(x float64) float64 {
	g.addPlatform(x+92, 172, 96, 32, tilePlatform)
	g.addEnemyOnSurface(enemyTurret, x+124, 172)
	g.addEnemyOnSurface(enemyGround, x+214, groundY)
	return 260
}

func (g *Game) chunkSkyBridge(x float64) float64 {
	g.addPlatform(x+68, 180, 64, 32, tilePlatform)
	g.addPlatform(x+156, 148, 96, 32, tilePlatform)
	g.addEnemy(enemyFlyer, x+168, 114)
	g.addEnemyOnSurface(enemyTurret, x+192, 148)
	return 290
}

func (g *Game) chunkMixedGauntlet(x float64) float64 {
	g.addPlatform(x+48, groundY-32, 32, 32, tileCrate)
	g.addHazard(x+102, groundY-18, 32, 18, tileSpike)
	g.addPlatform(x+156, 172, 64, 32, tilePlatform)
	g.addEnemyOnSurface(enemyTurret, x+174, 172)
	g.addEnemy(enemyFlyer, x+248, 104)
	g.addEnemyOnSurface(enemyGround, x+286, groundY)
	return 340
}

func (g *Game) addPlatform(x, y, w, h float64, tile tileKind) {
	g.platforms = append(g.platforms, platform{x: x, y: y, w: w, h: h, tile: tile})
}

func (g *Game) addHazard(x, y, w, h float64, tile tileKind) {
	g.hazards = append(g.hazards, hazard{x: x, y: y, w: w, h: h, tile: tile})
}

func (g *Game) addEnemyOnSurface(kind enemyKind, x, surfaceY float64) {
	g.addEnemy(kind, x, surfaceY-enemyHeight(kind))
}

func (g *Game) addEnemy(kind enemyKind, x, y float64) {
	e := enemy{
		id:      g.nextEnemyID,
		kind:    kind,
		element: enemyElement(kind),
		x:       x,
		y:       y,
		baseY:   y,
		phase:   g.rng.Float64() * math.Pi * 2,
	}
	g.nextEnemyID++

	switch kind {
	case enemyGround:
		e.w, e.h = 26, 20
		e.hp, e.maxHP = 3, 3
	case enemyTurret:
		e.w, e.h = 24, 28
		e.hp, e.maxHP = 4, 4
		e.shootTimer = 60 + g.rng.Intn(36)
	case enemyFlyer:
		e.w, e.h = 24, 20
		e.hp, e.maxHP = 3, 3
	}

	g.enemies = append(g.enemies, e)
}

func (g *Game) Draw(screen *ebiten.Image) {
	ox, oy := 0.0, 0.0
	if g.state == statePlaying {
		ox, oy = g.shakeOffset()
	}

	g.drawWorld(screen, ox, oy)

	switch g.state {
	case statePlaying:
		g.drawHUD(screen)
	case stateMenu:
		g.drawMenu(screen)
	case stateGameOver:
		g.drawGameOver(screen)
	}
}

func (g *Game) drawWorld(screen *ebiten.Image, ox, oy float64) {
	g.drawRepeating(screen, g.assets.backgroundFar, g.bgFarOffset+ox*0.15, oy*0.1)
	g.drawRepeating(screen, g.assets.backgroundMid, g.bgMidOffset+ox*0.35, oy*0.2)
	g.drawRepeating(screen, g.assets.backgroundNear, g.bgNearOffset+ox*0.55, oy*0.35)
	g.drawGround(screen, ox, oy)

	for _, p := range g.platforms {
		g.drawTiled(screen, g.assets.tiles[p.tile], p.x+ox, p.y+oy, p.w, p.h)
	}
	for _, h := range g.hazards {
		g.drawTiled(screen, g.assets.tiles[h.tile], h.x+ox, h.y+oy, h.w, h.h)
	}
	for _, e := range g.enemies {
		g.drawEnemy(screen, e, ox, oy)
	}
	for _, p := range g.projectiles {
		g.drawProjectile(screen, p, ox, oy)
	}
	g.drawPlayer(screen, ox, oy)
	for _, s := range g.swings {
		g.drawSwing(screen, s, ox, oy)
	}
}

func (g *Game) drawRepeating(screen *ebiten.Image, img *ebiten.Image, offsetX, offsetY float64) {
	w, _ := img.Size()
	x := math.Mod(offsetX, float64(w))
	if x > 0 {
		x -= float64(w)
	}
	for x < ScreenWidth {
		op := &ebiten.DrawImageOptions{}
		op.Filter = ebiten.FilterNearest
		op.GeoM.Translate(x, offsetY)
		screen.DrawImage(img, op)
		x += float64(w)
	}
}

func (g *Game) drawGround(screen *ebiten.Image, ox, oy float64) {
	ebitenutil.DrawRect(screen, 0, groundY+oy, ScreenWidth, ScreenHeight-groundY, color.RGBA{R: 40, G: 25, B: 25, A: 255})
	start := -math.Mod(g.distance, 32)
	for x := start - 32; x < ScreenWidth+32; x += 32 {
		op := &ebiten.DrawImageOptions{}
		op.Filter = ebiten.FilterNearest
		op.GeoM.Translate(x+ox, groundY+oy)
		screen.DrawImage(g.assets.tiles[tileGround], op)
	}
}

func (g *Game) drawTiled(screen *ebiten.Image, tile *ebiten.Image, x, y, w, h float64) {
	tw, th := tile.Size()
	for dx := 0.0; dx < w; dx += float64(tw) {
		op := &ebiten.DrawImageOptions{}
		op.Filter = ebiten.FilterNearest
		op.GeoM.Scale(1, h/float64(th))
		op.GeoM.Translate(x+dx, y)
		screen.DrawImage(tile, op)
	}
}

func (g *Game) drawEnemy(screen *ebiten.Image, e enemy, ox, oy float64) {
	if e.dead {
		return
	}
	frames := g.assets.enemyFrames[e.kind]
	frame := frames[(g.ticks/10)%len(frames)]
	if e.kind == enemyTurret {
		frame = frames[(g.ticks/18)%len(frames)]
	}

	op := &ebiten.DrawImageOptions{}
	op.Filter = ebiten.FilterNearest
	op.GeoM.Translate(e.x-4+ox, e.y-6+oy)
	if e.flash > 0 {
		op.ColorScale.Scale(1.35, 1.35, 1.35, 1)
	}
	if e.freeze > 0 {
		op.ColorScale.Scale(0.85, 1.1, 1.2, 1)
	}
	if e.stun > 0 {
		op.ColorScale.Scale(1.2, 1.15, 0.8, 1)
	}
	screen.DrawImage(frame, op)

	icon := g.assets.icons[e.element]
	iconOp := &ebiten.DrawImageOptions{}
	iconOp.Filter = ebiten.FilterNearest
	iconOp.GeoM.Translate(e.x+5+ox, e.y-14+oy)
	screen.DrawImage(icon, iconOp)

	if e.hp < e.maxHP {
		ebitenutil.DrawRect(screen, e.x+ox, e.y-6+oy, e.w, 3, color.RGBA{R: 35, G: 22, B: 28, A: 220})
		healthWidth := e.w * float64(e.hp) / float64(e.maxHP)
		ebitenutil.DrawRect(screen, e.x+ox, e.y-6+oy, healthWidth, 3, elementColors[e.element])
	}
}

func (g *Game) drawProjectile(screen *ebiten.Image, p projectile, ox, oy float64) {
	if p.dead {
		return
	}
	img := g.assets.enemyProjectiles[p.element]
	if p.fromPlayer {
		img = g.assets.playerProjectiles[p.element]
	}
	op := &ebiten.DrawImageOptions{}
	op.Filter = ebiten.FilterNearest
	sw, sh := img.Size()
	op.GeoM.Scale(p.w/float64(sw), p.h/float64(sh))
	op.GeoM.Translate(p.x+ox, p.y+oy)
	if !p.fromPlayer {
		op.ColorScale.Scale(0.92, 0.92, 0.92, 1)
	}
	screen.DrawImage(img, op)
}

func (g *Game) drawPlayer(screen *ebiten.Image, ox, oy float64) {
	frame := g.assets.playerFrames[0]
	if g.state == stateMenu {
		frame = g.assets.playerFrames[(g.ticks/12)%6]
	} else if !g.player.onGround {
		if g.player.vy < 0 {
			frame = g.assets.playerFrames[6]
		} else {
			frame = g.assets.playerFrames[7]
		}
	} else {
		frame = g.assets.playerFrames[(g.ticks/6)%6]
	}

	op := &ebiten.DrawImageOptions{}
	op.Filter = ebiten.FilterNearest
	op.GeoM.Translate(playerScreenX-8+ox, g.player.y-8+oy)
	if g.player.invuln > 0 && (g.player.invuln/4)%2 == 0 {
		op.ColorScale.ScaleAlpha(0.45)
	}
	screen.DrawImage(frame, op)

	if g.player.attackFlash > 0 {
		accent := elementColors[g.player.lastCast]
		glow := &ebiten.DrawImageOptions{}
		glow.Filter = ebiten.FilterNearest
		glow.GeoM.Translate(playerScreenX+g.player.w-2+ox, g.player.y+10+oy)
		glow.ColorScale.ScaleAlpha(float32(g.player.attackFlash) / 12)
		glow.ColorScale.Scale(ratio(accent.R), ratio(accent.G), ratio(accent.B), 1)
		screen.DrawImage(g.assets.icons[g.player.lastCast], glow)
	}
}

func (g *Game) drawSwing(screen *ebiten.Image, s swing, ox, oy float64) {
	if s.ttl <= 0 {
		return
	}
	img := g.assets.attackFX[s.element]
	r := s.rect(g.player)
	op := &ebiten.DrawImageOptions{}
	op.Filter = ebiten.FilterNearest
	sw, sh := img.Size()
	op.GeoM.Scale(r.W/float64(sw), r.H/float64(sh))
	op.GeoM.Translate(r.X+ox, r.Y+oy)
	op.ColorScale.ScaleAlpha(0.6 + float32(s.ttl)*0.05)
	screen.DrawImage(img, op)
}

func (g *Game) drawHUD(screen *ebiten.Image) {
	g.drawHUDPanel(screen, 8, 8, 92, 22, 210)
	for i := 0; i < maxHealth; i++ {
		x := 11 + float64(i*15)
		op := &ebiten.DrawImageOptions{}
		op.Filter = ebiten.FilterNearest
		op.GeoM.Scale(0.9, 0.9)
		op.GeoM.Translate(x, 10)
		if i >= g.player.hp {
			op.ColorScale.ScaleAlpha(0.2)
		}
		screen.DrawImage(g.assets.heart, op)
	}

	g.drawHUDPanel(screen, 300, 8, 172, 32, 210)
	g.drawSmallStatLine(screen, 308, 13, "Score", fmt.Sprintf("%04d", g.score), color.RGBA{R: 247, G: 233, B: 92, A: 255})
	g.drawSmallStatLine(screen, 390, 13, "Kills", fmt.Sprintf("%02d", g.kills), color.RGBA{R: 255, G: 105, B: 120, A: 255})
	g.drawSmallStatLine(screen, 308, 24, "Dist", fmt.Sprintf("%04d", int(g.distance/12)), color.RGBA{R: 126, G: 230, B: 255, A: 255})
	g.drawSmallStatLine(screen, 390, 24, "Best", fmt.Sprintf("%04d", g.bestScore), color.RGBA{R: 214, G: 226, B: 235, A: 255})

	if g.showAbilityHUD() {
		alpha := g.abilityHUDAlpha()
		y := 34 + (1-g.hudAbility)*(-12)
		g.drawHUDPanel(screen, 8, y, 92, 40, alpha)
		for i, el := range []element{fire, ice, thunder} {
			rowY := y + 4 + float64(i*12)
			g.drawAbilitySlot(screen, 12, rowY, 84, 8, el, alpha)
		}
	}
}

func (g *Game) drawGameOver(screen *ebiten.Image) {
	g.drawBackdrop(screen, 156)
	g.drawPanel(screen, 108, 76, 264, 118, color.RGBA{R: 255, G: 105, B: 120, A: 255})
	g.drawCenteredText(screen, "RUN ENDED", 108, 92, 264)
	g.drawCenteredText(screen, fmt.Sprintf("Score %d   Best %d", g.score, maxInt(g.bestScore, g.score)), 108, 116, 264)
	g.drawCenteredText(screen, "Enter or R to restart", 108, 142, 264)
	g.drawCenteredText(screen, "Esc to return to menu", 108, 158, 264)
	markerAccent := color.RGBA{R: 255, G: 141, B: 67, A: 255}
	ebitenutil.DrawRect(screen, 130+float64((g.ticks/20)%3), 144, 8, 8, markerAccent)
}

func (g *Game) drawMenu(screen *ebiten.Image) {
	g.drawBackdrop(screen, 140)
	g.drawPanel(screen, 98, 28, 284, 194, color.RGBA{R: 255, G: 141, B: 67, A: 255})

	leftIcon := &ebiten.DrawImageOptions{}
	leftIcon.Filter = ebiten.FilterNearest
	leftIcon.GeoM.Scale(1.4, 1.4)
	leftIcon.GeoM.Translate(124, 44)
	screen.DrawImage(g.assets.icons[fire], leftIcon)

	rightIcon := &ebiten.DrawImageOptions{}
	rightIcon.Filter = ebiten.FilterNearest
	rightIcon.GeoM.Scale(1.4, 1.4)
	rightIcon.GeoM.Translate(336, 44)
	screen.DrawImage(g.assets.icons[thunder], rightIcon)

	g.drawCenteredText(screen, "ELEMENT RUSH", 98, 50, 284)
	g.drawCenteredText(screen, "Elemental infinite runner", 98, 70, 284)
	g.drawCenteredText(screen, fmt.Sprintf("Best score %d", g.bestScore), 98, 88, 284)

	for i, label := range menuLabels {
		selected := i == g.menuIndex
		accent := color.RGBA{R: 255, G: 141, B: 67, A: 255}
		if i == menuQuit {
			accent = color.RGBA{R: 255, G: 105, B: 120, A: 255}
		}
		g.drawMenuOption(screen, 138, float64(108+i*36), 204, 28, label, selected, accent)
	}

	g.drawCenteredText(screen, "Move  W/S or Up/Down", 98, 186, 284)
	g.drawCenteredText(screen, "Enter/Space select   Esc quit", 98, 200, 284)
}

func (g *Game) drawMenuOption(screen *ebiten.Image, x, y, w, h float64, label string, selected bool, accent color.RGBA) {
	fill := color.RGBA{R: 20, G: 24, B: 34, A: 220}
	if selected {
		fill = color.RGBA{R: 34, G: 40, B: 54, A: 235}
	}
	ebitenutil.DrawRect(screen, x, y, w, h, fill)
	g.drawRectOutline(screen, x, y, w, h, color.RGBA{R: 214, G: 226, B: 235, A: 230})
	if selected {
		ebitenutil.DrawRect(screen, x, y, w, 3, accent)
		ebitenutil.DrawRect(screen, x+7+float64((g.ticks/18)%2), y+7, 4, h-14, accent)
	}
	g.drawCenteredText(screen, label, int(x), int(y)+7, int(w))
}

func (g *Game) drawAbilitySlot(screen *ebiten.Image, x, y, w, h float64, el element, alpha uint8) {
	accent := elementColors[el]
	selected := g.player.lastCast == el
	fill := color.RGBA{R: 19, G: 23, B: 31, A: 0}
	if selected {
		fill = color.RGBA{R: 28, G: 33, B: 42, A: uint8(minInt(int(alpha), 96))}
	}
	ebitenutil.DrawRect(screen, x, y, w, h, fill)
	if selected {
		ebitenutil.DrawRect(screen, x-2, y, 2, h, color.RGBA{R: accent.R, G: accent.G, B: accent.B, A: alpha})
	}
	if selected {
		g.drawRectOutline(screen, x, y, w, h, color.RGBA{R: accent.R, G: accent.G, B: accent.B, A: uint8(minInt(int(alpha), 170))})
	}

	icon := &ebiten.DrawImageOptions{}
	icon.Filter = ebiten.FilterNearest
	icon.GeoM.Scale(0.7, 0.7)
	icon.GeoM.Translate(x, y-2)
	icon.ColorScale.ScaleAlpha(float32(alpha) / 255)
	if g.player.cooldowns[int(el)] > 0 {
		icon.ColorScale.ScaleAlpha(0.55)
	}
	screen.DrawImage(g.assets.icons[el], icon)

	cdMax := attackCooldown(el)
	cdNow := g.player.cooldowns[int(el)]
	barX := x + 15
	barY := y + 2
	barW := w - 19
	ebitenutil.DrawRect(screen, barX, barY, barW, 4, color.RGBA{R: 13, G: 17, B: 22, A: alpha})
	fillW := barW
	if cdNow > 0 {
		fillW = barW * (1 - float64(cdNow)/float64(cdMax))
	}
	ebitenutil.DrawRect(screen, barX, barY, fillW, 4, color.RGBA{R: accent.R, G: accent.G, B: accent.B, A: alpha})
}

func (g *Game) drawBackdrop(screen *ebiten.Image, alpha uint8) {
	ebitenutil.DrawRect(screen, 0, 0, ScreenWidth, ScreenHeight, color.RGBA{R: 8, G: 10, B: 16, A: alpha})
}

func (g *Game) drawPanel(screen *ebiten.Image, x, y, w, h float64, accent color.RGBA) {
	ebitenutil.DrawRect(screen, x+3, y+3, w, h, color.RGBA{R: 0, G: 0, B: 0, A: 70})
	ebitenutil.DrawRect(screen, x, y, w, h, color.RGBA{R: 18, G: 22, B: 30, A: 230})
	g.drawRectOutline(screen, x, y, w, h, color.RGBA{R: 214, G: 226, B: 235, A: 220})
	ebitenutil.DrawRect(screen, x+1, y+1, w-2, 3, accent)
	ebitenutil.DrawRect(screen, x+1, y+h-4, w-2, 1, color.RGBA{R: 255, G: 255, B: 255, A: 30})
}

func (g *Game) drawHUDPanel(screen *ebiten.Image, x, y, w, h float64, alpha uint8) {
	ebitenutil.DrawRect(screen, x+2, y+2, w, h, color.RGBA{R: 0, G: 0, B: 0, A: uint8(minInt(int(alpha/4), 44))})
	ebitenutil.DrawRect(screen, x, y, w, h, color.RGBA{R: 15, G: 19, B: 27, A: alpha})
	g.drawRectOutline(screen, x, y, w, h, color.RGBA{R: 220, G: 228, B: 236, A: uint8(minInt(int(alpha)+8, 220))})
}

func (g *Game) drawSmallStatLine(screen *ebiten.Image, x, y int, label, value string, accent color.RGBA) {
	ebitenutil.DrawRect(screen, float64(x), float64(y+2), 4, 4, accent)
	labelX := x + 8
	valueX := labelX + smallTextWidth(label) + 8
	g.drawSmallText(screen, label, labelX, y, color.RGBA{R: accent.R, G: accent.G, B: accent.B, A: 255})
	g.drawSmallText(screen, value, valueX, y, color.RGBA{R: 245, G: 247, B: 250, A: 255})
}

func (g *Game) drawRectOutline(screen *ebiten.Image, x, y, w, h float64, c color.RGBA) {
	ebitenutil.DrawRect(screen, x, y, w, 1, c)
	ebitenutil.DrawRect(screen, x, y+h-1, w, 1, c)
	ebitenutil.DrawRect(screen, x, y, 1, h, c)
	ebitenutil.DrawRect(screen, x+w-1, y, 1, h, c)
}

func (g *Game) drawCenteredText(screen *ebiten.Image, text string, x, y, w int) {
	g.drawText(screen, text, x+(w-textWidth(text))/2, y)
}

func (g *Game) drawSmallText(screen *ebiten.Image, text string, x, y int, clr color.Color) {
	metrics := hudSmallFont.Metrics()
	ebitext.Draw(screen, text, hudSmallFont, x, y+metrics.Ascent.Ceil(), clr)
}

func (g *Game) drawText(screen *ebiten.Image, text string, x, y int) {
	ebitenutil.DebugPrintAt(screen, text, x+1, y+1)
	ebitenutil.DebugPrintAt(screen, text, x, y)
}

func (g *Game) shakeOffset() (float64, float64) {
	if g.shake <= 0 {
		return 0, 0
	}
	rangeX := float64(g.shake) * 0.5
	return g.rng.Float64()*rangeX - rangeX/2, g.rng.Float64()*rangeX - rangeX/2
}

func (g *Game) updateAbilityHUDAnimation() {
	target := 0.0
	if g.hudFocus > 0 || g.hasActiveCooldown() {
		target = 1.0
	}
	if g.hudAbility < target {
		g.hudAbility = math.Min(target, g.hudAbility+0.14)
		return
	}
	g.hudAbility = math.Max(target, g.hudAbility-0.1)
}

func (g *Game) showAbilityHUD() bool {
	return g.hudAbility > 0.02
}

func (g *Game) abilityHUDAlpha() uint8 {
	return uint8(220 * g.hudAbility)
}

func (g *Game) hasActiveCooldown() bool {
	for _, cd := range g.player.cooldowns {
		if cd > 0 {
			return true
		}
	}
	return false
}

func attackCooldown(el element) int {
	switch el {
	case fire:
		return 20
	case ice:
		return 30
	default:
		return 15
	}
}

func attackDamage(el element) int {
	switch el {
	case fire:
		return 2
	case ice:
		return 2
	default:
		return 1
	}
}

func attackVelocity(el element) float64 {
	switch el {
	case fire:
		return 7.0
	case ice:
		return 5.6
	default:
		return 8.6
	}
}

func attackProjectileTTL(el element) int {
	switch el {
	case fire:
		return 52
	case ice:
		return 68
	default:
		return 42
	}
}

func attackPierce(el element) int {
	if el == thunder {
		return 2
	}
	return 1
}

func attackProjectileSize(el element) (float64, float64) {
	switch el {
	case fire:
		return 16, 16
	case ice:
		return 20, 14
	default:
		return 18, 14
	}
}

func attackProjectileSpawnY(el element) float64 {
	switch el {
	case fire:
		return 12
	case ice:
		return 14
	default:
		return 10
	}
}

func attackSwing(el element) (float64, float64, float64, float64) {
	switch el {
	case fire:
		return 36, 28, 20, 8
	case ice:
		return 42, 30, 16, 6
	default:
		return 34, 26, 22, 7
	}
}

func attackSwingTTL(el element) int {
	switch el {
	case fire:
		return 8
	case ice:
		return 10
	default:
		return 6
	}
}

func enemyElement(kind enemyKind) element {
	switch kind {
	case enemyGround:
		return ice
	case enemyTurret:
		return fire
	default:
		return thunder
	}
}

func enemyHeight(kind enemyKind) float64 {
	switch kind {
	case enemyGround:
		return 20
	case enemyTurret:
		return 28
	default:
		return 20
	}
}

func isStrongAgainst(attacker, target element) bool {
	return (attacker == fire && target == ice) || (attacker == ice && target == thunder) || (attacker == thunder && target == fire)
}

func attackKey(el element) string {
	switch el {
	case fire:
		return "J"
	case ice:
		return "K"
	default:
		return "L"
	}
}

func wrapLayer(value float64) float64 {
	for value <= -ScreenWidth {
		value += ScreenWidth
	}
	for value > 0 {
		value -= ScreenWidth
	}
	return value
}

func ratio(v uint8) float32 {
	return float32(v) / 255
}

func textWidth(text string) int {
	return len(text) * 6
}

func smallTextWidth(text string) int {
	return ebitext.BoundString(hudSmallFont, text).Dx()
}

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}
