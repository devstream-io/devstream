package trello

var IssuesBuilder = `
name: Issues Builder
on:
  issues:
    types: [opened, reopened, edited, closed]
  issue_comment:
    types: [created, edited, deleted]
jobs:
  issue_comment_integration:
    if: ${{ github.event_name == 'issue_comment' }}
    runs-on: ubuntu-latest
    name: Integrate Issue Comment
    env:
      ACTION_VERSION: CyanRYi/trello-github-integration@v3.0.0
      TRELLO_API_KEY: ${{ secrets.TRELLO_API_KEY }}
      TRELLO_API_TOKEN: ${{ secrets.TRELLO_TOKEN }}
      TRELLO_BOARD_ID: ${{ secrets.TRELLO_BOARD_ID }}
      TRELLO_TODO_LIST_ID: ${{ secrets.TRELLO_TODO_LIST_ID }}
      TRELLO_DONE_LIST_ID: ${{ secrets.TRELLO_DONE_LIST_ID }}
      TRELLO_MEMBER_MAP: ${{ secrets.TRELLO_MEMBER_MAP }}
    steps:
      - name: Add Comment
        id: comment-create
        if: ${{ github.event.action == 'created' }}
        uses: CyanRYi/trello-github-integration@v3.0.0
        with:
          trello-action: add_comment
      - name: Edit Comment
        id: comment-edit
        if: ${{ github.event.action == 'edited' }}
        uses: CyanRYi/trello-github-integration@v3.0.0
        with:
          trello-action: edit_comment
      - name: Delete Comment
        id: comment-delete
        if: ${{ github.event.action == 'deleted' }}
        uses: CyanRYi/trello-github-integration@v3.0.0
        with:
          trello-action: delete_comment
  issue_integration:
    if: ${{ github.event_name == 'issues' }}
    runs-on: ubuntu-latest
    name: Integrate Issue
    env:
      TRELLO_API_KEY: ${{ secrets.TRELLO_API_KEY }}
      TRELLO_API_TOKEN: ${{ secrets.TRELLO_TOKEN }}
      TRELLO_BOARD_ID: ${{ secrets.TRELLO_BOARD_ID }}
      TRELLO_TODO_LIST_ID: ${{ secrets.TRELLO_TODO_LIST_ID }}
      TRELLO_DONE_LIST_ID: ${{ secrets.TRELLO_DONE_LIST_ID }}
      TRELLO_MEMBER_MAP: ${{ secrets.TRELLO_MEMBER_MAP }}
    steps:
      - name: Create Card
        id: card-create
        if: ${{ github.event.action == 'opened' }}
        uses: CyanRYi/trello-github-integration@v3.0.0
        with:
          trello-action: create_card
      - name: Edit Card
        id: card-edit
        if: ${{ github.event.action == 'edited' }}
        uses: CyanRYi/trello-github-integration@v3.0.0
        with:
          trello-action: edit_card
      - name: Close Card
        id: card-move-to-done
        if: ${{ github.event.action == 'closed' }}
        uses: CyanRYi/trello-github-integration@v3.0.0
        with:
          trello-action: move_card_to_done
      - name: Reopen Card
        id: card-move-to-backlog
        if: ${{ github.event.action == 'reopened' }}
        uses: CyanRYi/trello-github-integration@v3.0.0
        with:
          trello-action: move_card_to_todo
`
