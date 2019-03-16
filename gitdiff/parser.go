	"os"
	fileHeaderPrefix = "diff --git "
	// TODO(bkeyes): parse header line for filename
	// necessary to get the filename for mode changes or add/rm empty files

	for {
		line, err := p.PeekLine()
		if err != nil {
			return err
		}

		more, err := parseGitHeaderLine(f, line)
		if err != nil {
			return p.Errorf("header: %v", err)
		}
		if !more {
			break
		}
		p.Line()
	}

	return nil

func parseGitHeaderLine(f *File, line string) (more bool, err error) {
	match := func(s string) bool {
		if strings.HasPrefix(line, s) {
			// TODO(bkeyes): strip final line separator too
			line = line[len(s):]
			return true
		}
		return false
	}

	switch {
	case match(fragmentHeaderPrefix):
		// start of a fragment indicates the end of the header
		return false, nil

	case match(oldFilePrefix):

	case match(newFilePrefix):

	case match("old mode "):
		if f.OldMode, err = parseModeLine(line); err != nil {
			return false, err
		}

	case match("new mode "):
		if f.NewMode, err = parseModeLine(line); err != nil {
			return false, err
		}

	case match("deleted file mode "):
		// TODO(bkeyes): maybe set old name from default?
		f.IsDelete = true
		if f.OldMode, err = parseModeLine(line); err != nil {
			return false, err
		}

	case match("new file mode "):
		f.IsNew = true
		if f.NewMode, err = parseModeLine(line); err != nil {
			return false, err
		}

	case match("copy from "):
		f.IsCopy = true
		// TODO(bkeyes): set old name

	case match("copy to "):
		f.IsCopy = true
		// TODO(bkeyes): set new name

	case match("rename old "):
		f.IsRename = true
		// TODO(bkeyes): set old name

	case match("rename new "):
		f.IsRename = true
		// TODO(bkeyes): set new name

	case match("rename from "):
		f.IsRename = true
		// TODO(bkeyes): set old name

	case match("rename to "):
		f.IsRename = true
		// TODO(bkeyes): set new name

	case match("similarity index "):
		f.Score = parseScoreLine(line)

	case match("dissimilarity index "):
		f.Score = parseScoreLine(line)

	case match("index "):

	default:
		// unknown line also indicates the end of the header
		return false, nil
	}

	return true, nil
}

func parseModeLine(s string) (os.FileMode, error) {
	s = strings.TrimSuffix(s, "\n")

	mode, err := strconv.ParseInt(s, 8, 32)
	if err != nil {
		nerr := err.(*strconv.NumError)
		return os.FileMode(0), fmt.Errorf("invalid mode line: %v", nerr.Err)
	}

	return os.FileMode(mode), nil
}

func parseScoreLine(s string) int {
	s = strings.TrimSuffix(s, "\n")

	// gitdiff_similarity / gitdiff_dissimilarity ignore invalid scores
	score, _ := strconv.ParseInt(s, 10, 32)
	if score <= 100 {
		return int(score)
	}
	return 0
}