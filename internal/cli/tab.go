package cli

import (
	"fmt"
	"strings"
	"time"

	"github.com/freecode/freecode/internal/config"
	"github.com/freecode/freecode/internal/session"
	"github.com/google/uuid"
	"github.com/spf13/cobra"
)

var tabCmd = &cobra.Command{
	Use:   "tab",
	Short: "Manage tabs",
	Long:  `Create, close, and manage session tabs.`,
}

var (
	tabForce bool
)

func init() {
	tabCmd.AddCommand(tabNewCmd)
	tabCmd.AddCommand(tabCloseCmd)
	tabCmd.AddCommand(tabListCmd)
	tabCmd.AddCommand(tabMoveCmd)
	tabCmd.AddCommand(tabRenameCmd)

	tabCloseCmd.Flags().BoolVarP(&tabForce, "force", "f", false, "Skip confirmation prompt")
}

var tabNewCmd = &cobra.Command{
	Use:   "new [name]",
	Short: "Create a new tab",
	RunE:  runTabNew,
}

var tabCloseCmd = &cobra.Command{
	Use:   "close [tab-id]",
	Short: "Close a tab",
	RunE:  runTabClose,
}

var tabListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all tabs",
	RunE:  runTabList,
}

var tabMoveCmd = &cobra.Command{
	Use:   "move [session-id] [tab-id]",
	Short: "Move session to tab",
	RunE:  runTabMove,
}

var tabRenameCmd = &cobra.Command{
	Use:   "rename [tab-id] [name]",
	Short: "Rename a tab",
	RunE:  runTabRename,
}

func getTabStore() *session.Store {
	cfg := config.DefaultConfig()
	return session.NewStore(cfg.Session.Dir)
}

func runTabNew(cmd *cobra.Command, args []string) error {
	store := getTabStore()

	name := "New Tab"
	if len(args) > 0 {
		name = strings.Join(args, " ")
	}

	tab := &session.Tab{
		ID:        uuid.New().String(),
		Name:      name,
		CreatedAt: time.Now(),
		Sessions:  make([]string, 0),
	}

	if err := store.SaveTab(tab); err != nil {
		return fmt.Errorf("failed to create tab: %w", err)
	}

	fmt.Printf("Created tab %s: %s\n", tab.ID, tab.Name)
	return nil
}

func runTabClose(cmd *cobra.Command, args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("tab id is required")
	}

	tabID := args[0]
	store := getTabStore()

	tabs, err := store.LoadTabs()
	if err != nil {
		return fmt.Errorf("failed to load tabs: %w", err)
	}

	var tab *session.Tab
	for _, t := range tabs {
		if t.ID == tabID {
			tab = t
			break
		}
	}

	if tab == nil {
		return fmt.Errorf("tab not found: %s", tabID)
	}

	if !tabForce {
		fmt.Printf("Are you sure you want to close tab %s (%s)? This will not delete sessions. (y/N): ", tab.ID, tab.Name)
		var confirm string
		fmt.Scanln(&confirm)
		if strings.ToLower(confirm) != "y" {
			fmt.Println("Cancelled")
			return nil
		}
	}

	if err := store.DeleteTab(tabID); err != nil {
		return fmt.Errorf("failed to close tab: %w", err)
	}

	fmt.Printf("Closed tab: %s\n", tab.Name)
	return nil
}

func runTabList(cmd *cobra.Command, args []string) error {
	store := getTabStore()
	tabStore := session.NewStore(config.DefaultConfig().Session.Dir)

	tabs, err := store.LoadTabs()
	if err != nil {
		return fmt.Errorf("failed to load tabs: %w", err)
	}

	if len(tabs) == 0 {
		fmt.Println("No tabs found.")
		return nil
	}

	fmt.Println("Tabs:")
	fmt.Println("")
	fmt.Printf("%-38s %-20s %s\n", "ID", "Created", "Name")
	fmt.Println(strings.Repeat("-", 80))

	for _, tab := range tabs {
		id := tab.ID
		if len(id) > 36 {
			id = id[:36]
		}
		created := tab.CreatedAt.Format("2006-01-02 15:04")
		name := tab.Name
		if name == "" {
			name = "(unnamed)"
		}
		fmt.Printf("%-38s %-20s %s\n", id, created, name)

		sessions, err := tabStore.ListSessions()
		if err == nil {
			var tabSessions []string
			for _, s := range sessions {
				if s.TabID == tab.ID {
					tabSessions = append(tabSessions, s.ID)
				}
			}
			if len(tabSessions) > 0 {
				fmt.Printf("  Sessions: %d\n", len(tabSessions))
			} else {
				fmt.Println("  Sessions: 0")
			}
		}
	}

	fmt.Println("")
	fmt.Printf("Total: %d tab(s)\n", len(tabs))

	return nil
}

func runTabMove(cmd *cobra.Command, args []string) error {
	if len(args) < 2 {
		return fmt.Errorf("session-id and tab-id are required")
	}

	sessionID := args[0]
	tabID := args[1]

	store := getTabStore()

	tabs, err := store.LoadTabs()
	if err != nil {
		return fmt.Errorf("failed to load tabs: %w", err)
	}

	var tab *session.Tab
	for _, t := range tabs {
		if t.ID == tabID {
			tab = t
			break
		}
	}

	if tab == nil {
		return fmt.Errorf("tab not found: %s", tabID)
	}

	sessions, err := store.ListSessions()
	if err != nil {
		return fmt.Errorf("failed to load sessions: %w", err)
	}

	var sess *session.Session
	for _, s := range sessions {
		if s.ID == sessionID {
			sess = s
			break
		}
	}

	if sess == nil {
		return fmt.Errorf("session not found: %s", sessionID)
	}

	oldTabID := sess.TabID
	sess.TabID = tabID

	if err := store.SaveSession(sess); err != nil {
		return fmt.Errorf("failed to update session: %w", err)
	}

	tab.Sessions = append(tab.Sessions, sessionID)
	tab.ActiveSession = sessionID

	if err := store.SaveTab(tab); err != nil {
		return fmt.Errorf("failed to update tab: %w", err)
	}

	if oldTabID != "" {
		for _, t := range tabs {
			if t.ID == oldTabID {
				for i, sid := range t.Sessions {
					if sid == sessionID {
						t.Sessions = append(t.Sessions[:i], t.Sessions[i+1:]...)
						store.SaveTab(t)
						break
					}
				}
			}
		}
	}

	fmt.Printf("Moved session %s to tab %s\n", sessionID, tabID)
	return nil
}

func runTabRename(cmd *cobra.Command, args []string) error {
	if len(args) < 2 {
		return fmt.Errorf("tab-id and name are required")
	}

	tabID := args[0]
	name := strings.Join(args[1:], " ")

	store := getTabStore()

	tabs, err := store.LoadTabs()
	if err != nil {
		return fmt.Errorf("failed to load tabs: %w", err)
	}

	var tab *session.Tab
	for _, t := range tabs {
		if t.ID == tabID {
			tab = t
			break
		}
	}

	if tab == nil {
		return fmt.Errorf("tab not found: %s", tabID)
	}

	tab.Name = name

	if err := store.SaveTab(tab); err != nil {
		return fmt.Errorf("failed to rename tab: %w", err)
	}

	fmt.Printf("Renamed tab to: %s\n", name)
	return nil
}