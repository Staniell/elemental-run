package game

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type saveData struct {
	BestScore int `json:"best_score"`
}

func defaultSavePath() string {
	if dir, err := os.UserConfigDir(); err == nil && dir != "" {
		return filepath.Join(dir, "ElementRush", "save.json")
	}
	return "element-rush-save.json"
}

func readSaveData(path string) (saveData, error) {
	var data saveData
	bs, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return data, nil
		}
		return data, err
	}
	if err := json.Unmarshal(bs, &data); err != nil {
		return saveData{}, err
	}
	return data, nil
}

func writeSaveData(path string, data saveData) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	bs, err := json.Marshal(data)
	if err != nil {
		return err
	}
	return os.WriteFile(path, bs, 0o644)
}

func (g *Game) loadProgress() {
	if g.savePath == "" {
		g.savePath = defaultSavePath()
	}
	data, err := readSaveData(g.savePath)
	if err != nil {
		return
	}
	if data.BestScore > g.bestScore {
		g.bestScore = data.BestScore
	}
}

func (g *Game) setBestScore(score int) {
	if score <= g.bestScore {
		return
	}
	g.bestScore = score
	if g.savePath == "" {
		g.savePath = defaultSavePath()
	}
	_ = writeSaveData(g.savePath, saveData{BestScore: g.bestScore})
}
