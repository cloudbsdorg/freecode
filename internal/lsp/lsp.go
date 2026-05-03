package lsp

import (
	"context"
	"encoding/json"
)

type Client struct {
	conn   json.RawMessage
	server *Server
}

type Server struct {
	stdin  chan json.RawMessage
	stdout chan json.RawMessage
}

type Diagnostic struct {
	Range    Range  `json:"range"`
	Severity int    `json:"severity"`
	Message  string `json:"message"`
}

type Range struct {
	Start Position `json:"start"`
	End   Position `json:"end"`
}

type Position struct {
	Line      uint32 `json:"line"`
	Character uint32 `json:"character"`
}

func NewClient() *Client {
	return &Client{}
}

func (c *Client) Connect(ctx context.Context, server string) error {
	return nil
}

func (c *Client) Initialize(ctx context.Context) error {
	return nil
}

func (c *Client) Shutdown(ctx context.Context) error {
	return nil
}

func (c *Client) DidOpen(ctx context.Context, uri string, text string) error {
	return nil
}

func (c *Client) DidChange(ctx context.Context, uri string, text string) error {
	return nil
}

func (c *Client) Completion(ctx context.Context, uri string, pos Position) ([]CompletionItem, error) {
	return nil, nil
}

type CompletionItem struct {
	Label         string `json:"label"`
	InsertText    string `json:"insertText"`
	Kind          int    `json:"kind"`
	Documentation string `json:"documentation"`
}

func NewServer() *Server {
	return &Server{
		stdin:  make(chan json.RawMessage),
		stdout: make(chan json.RawMessage),
	}
}

func (s *Server) Start(ctx context.Context) error {
	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	return nil
}

func (s *Server) Send(ctx context.Context, msg json.RawMessage) error {
	return nil
}

func (s *Server) Receive(ctx context.Context) (json.RawMessage, error) {
	return nil, nil
}
