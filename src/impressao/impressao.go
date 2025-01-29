package impressao

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"

	"golang.org/x/text/encoding/charmap"
)

// Funções de formatação
func SetBold(enable bool) string {
	if enable {
		return "\x1b\x45\x01" // Ativar negrito
	}
	return "\x1b\x45\x00" // Desativar negrito
}

func SetFontSize(width, height int) string {
	return fmt.Sprintf("\x1d\x21%c", byte((width-1)<<4|height-1))
}

func ResetPrinter() string {
	return "\x1b\x40" // Reset da impressora
}

func SetCharacterSetCP850() string {
	return "\x1b\x74\x02" // Configura o conjunto de caracteres para CP850
}

func CutPaper() string {
	return "\x1d\x56\x42" // Comando para corte parcial
}

// Função para envio de dados para impressora
func PrintToPrinter(printerPath, text string) error {
	// Converter texto para CP850
	encoder := charmap.CodePage850.NewEncoder()
	encodedText, err := encoder.String(text)
	if err != nil {
		return fmt.Errorf("erro ao codificar texto: %v", err)
	}

	// Criar arquivo temporário para o spool de impressão
	tempFile := "temp_print.txt"
	err = os.WriteFile(tempFile, []byte(encodedText), 0644)
	if err != nil {
		return fmt.Errorf("erro ao criar arquivo temporário: %v", err)
	}

	// Enviar o arquivo para a impressora
	cmd := exec.Command("cmd", "/C", "copy", tempFile, printerPath)
	var output bytes.Buffer
	cmd.Stdout = &output
	cmd.Stderr = &output
	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("erro ao enviar para a impressora: %v\nSaída: %s", err, output.String())
	}

	// Remover o arquivo temporário
	err = os.Remove(tempFile)
	if err != nil {
		log.Printf("Erro ao remover arquivo temporário: %v", err)
	}

	return nil
}
