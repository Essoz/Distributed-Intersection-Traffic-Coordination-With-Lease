#!/usr/bin/env bash

session="ece445-demo"

# create a new tmux session called "ece445-demo", delete the session if it already exists
tmux list-panes -s -F '#{session_name}:#{window_index}.#{pane_index}' | while read pane; do
  tmux send-keys -t "$pane" C-c
done
tmux kill-session -t $session 2>/dev/null || true