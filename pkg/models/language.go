package models

import protocol "v/protocol/go_protocol"

type Language struct {
	Code string
	Label string
	TranscriberCode string
	TTSModel string
}

func NewLanguageModel(lang protocol.Language) (t Language){
	switch lang {
	case *protocol.Language_ARABIC.Enum():
		t = Language{
			Code:             "ar-AR",
			Label:            "Arabic",
			TranscriberCode:  "ar-AR",
			TTSModel: "ar-AR-Wavenet-D",
		}
	case *protocol.Language_FRENCH.Enum():
		t = Language{
			Code:             "fr-FR",
			Label:            "Français",
			TranscriberCode:  "fr-FR",
			TTSModel: "fr-FR-Wavenet-B",
		}
	case *protocol.Language_TURKISH.Enum():
		t = Language{
			Code:             "tr-TR",
			Label:            "Türkçe",
			TranscriberCode:  "tr-TR",
			TTSModel: "tr-TR-Wavenet-B",
		}
	case *protocol.Language_ENGLISH.Enum():
		fallthrough;
	default:
		t = Language{
			Code:             "en-US",
			Label:            "English",
			TranscriberCode:  "en-US",
			TTSModel: "en-US-Wavenet-D",
		}
	}
	
	return t
}
