package commands

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/ledongthuc/pdf"
	"github.com/reddy-santhu/study-bot/db"
)

func HandleUploadPDF(s *discordgo.Session, m *discordgo.MessageCreate) {
	if len(m.Attachments) > 0 {
		attachment := m.Attachments[0]
		if strings.HasSuffix(attachment.Filename, ".pdf") {
			err := DownloadAndSavePDF(attachment.URL, attachment.Filename)
			if err != nil {
				log.Printf("Error downloading and saving PDF: %v", err)
				s.ChannelMessageSend(m.ChannelID, "Error downloading and saving the PDF. Please try again.")
				return
			}
			text, err := extractTextFromPDF(attachment.Filename)
			if err != nil {
				log.Printf("Error with Text extraction: %v", err)
				s.ChannelMessageSend(m.ChannelID, "Error during text extraction please try again.")
				return
			}
			err = db.LogPDFData(m.Author.ID, attachment.Filename, text)
			if err != nil {
				log.Printf("Error logging to MongoDB: %v", err)
				s.ChannelMessageSend(m.ChannelID, "Issue Logging PDF to MongoDB")
				return
			}

			s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("PDF '%s' downloaded and logged in MongoDB!", attachment.Filename))

		} else {
			s.ChannelMessageSend(m.ChannelID, "Only PDF files are supported.")
		}
	} else {
		s.ChannelMessageSend(m.ChannelID, "Please attach a PDF file.")
	}
}

func DownloadAndSavePDF(url, filename string) error {
	response, err := http.Get(url)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		return fmt.Errorf("received non-200 status code: %d", response.StatusCode)
	}

	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, response.Body)
	if err != nil {
		return err
	}

	return nil
}
func extractTextFromPDF(pdfPath string) (string, error) {
	f, err := os.Open(pdfPath)
	if err != nil {
		return "", fmt.Errorf("error with opening pdf: %w", err)
	}
	defer f.Close()
	fi, err := f.Stat()
	if err != nil {
		return "", fmt.Errorf("error with reading PDF to extract content: %w", err)
	}
	r, err := pdf.NewReader(f, fi.Size())
	if err != nil {
		return "", fmt.Errorf("error with PDF reader: %w", err)
	}
	var textBuilder strings.Builder
	textReader, err := r.GetPlainText()
	_, err = io.Copy(&textBuilder, textReader)
	if err != nil {
		return "", fmt.Errorf("error extracting text from PDF: %w", err)
	}
	return textBuilder.String(), nil
}
func HandleViewPDF(s *discordgo.Session, m *discordgo.MessageCreate) {
	userID := m.Author.ID
	pdfs, err := db.GetPDFsByUser(userID)
	if err != nil {
		log.Printf("Error getting PDFs for user %s: %v", userID, err)
		s.ChannelMessageSend(m.ChannelID, "Error occurred while retrieving your PDFs. Please try again.")
		return
	}

	if len(pdfs) == 0 {
		s.ChannelMessageSend(m.ChannelID, "You have not uploaded any PDFs yet. Use /pdf to upload.")
		return
	}

	var message strings.Builder
	message.WriteString("Your Uploaded PDFs:\n")
	for i, pdfData := range pdfs {
		message.WriteString(fmt.Sprintf("%d. %s\n", i+1, pdfData.Filename)) // Format as a numbered list
	}

	s.ChannelMessageSend(m.ChannelID, message.String())
}
