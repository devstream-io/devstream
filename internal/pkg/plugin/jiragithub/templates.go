package jiragithub

var jiraIssuesBuilder = `name: jira-github-integ Builder 
on:
  issues:
    types: [opened, reopened, edited, closed]
  issue_comment:
    types: [created, edited, deleted]

jobs:
  github-jira-issue:
    name: Transition Issue
    runs-on: ubuntu-latest
    steps:
    - name: Login
      uses: atlassian/gajira-login@master
      env:
        JIRA_BASE_URL: [[.JiraBaseUrl]]
        JIRA_USER_EMAIL: [[.JiraUserEmail]]
        JIRA_API_TOKEN: ${{ secrets.JIRA_API_TOKEN }}

    - name: Create new Jira issue
      id: create
      if: ${{ github.event.action == 'opened' }}
      uses: atlassian/gajira-create@master
      with:
        project: [[.JiraProjectKey]]
        issuetype: Task
        summary: ${{ github.event.issue.title }}
        description: ${{ github.event.issue.body }}
    
    - name: Rename github issue
      if: ${{ github.event.action == 'opened' }}
      uses: actions-cool/issues-helper@v3
      with:
        actions: 'update-issue'
        token: ${{ secrets.GH_TOKEN }}
        issue-number: ${{ github.event.issue.number }}
        state: 'open'
        title: "[${{ steps.create.outputs.issue }}] ${{ github.event.issue.title }}"
        update-mode: 'replace'
        emoji: '+1'
      
    - name: Find Jira Issue Key
      id: find
      if: ${{ github.event.action != 'opened' }}
      uses: atlassian/gajira-find-issue-key@master
      with:
        string: "${{ github.event.issue.title }}"
        
    - name: Transition issue to In Progress
      id: transition_inprogress
      if: ${{ github.event.action == 'edited' }}
      uses: atlassian/gajira-transition@master
      with:
        issue: ${{ steps.find.outputs.issue }}
        transition: "In Progress"

    - name: Transition issue to Done
      id: transition_done
      if: ${{ github.event.action == 'closed' }}
      uses: atlassian/gajira-transition@master
      with:
        issue: ${{ steps.find.outputs.issue }}
        transition: "Done"

    - name: Transition issue to TO DO
      id: transition_todo
      if: ${{ github.event.action == 'reopened' }}
      uses: atlassian/gajira-transition@master
      with:
        issue: ${{ steps.find.outputs.issue }}
        transition: "To Do"
    
  issue_comment_integration:
    if: ${{ github.event_name == 'issue_comment' }}
    runs-on: ubuntu-latest
    name: Integrate Issue Comment
    steps:
    - name: Login
      uses: atlassian/gajira-login@master
      env:
        JIRA_BASE_URL: [[.JiraBaseUrl]]
        JIRA_USER_EMAIL: [[.JiraUserEmail]]
        JIRA_API_TOKEN: ${{ secrets.JIRA_API_TOKEN }}
    - name: Find Issue Key
      id: find
      uses: atlassian/gajira-find-issue-key@master
      with:
        string: "${{ github.event.issue.title }}"        
    - name: Create issue
      id: create_issue
      if: ${{ github.event.action == 'created' }}
      uses: atlassian/gajira-comment@master
      with:
        issue: ${{ steps.find.outputs.issue }}
        comment: ${{ github.event.comment.body }}

    - name: Edit issue
      id: edit_issue
      if: ${{ github.event.action == 'edited' }}
      uses: atlassian/gajira-comment@master
      with:
        issue: ${{ steps.find.outputs.issue }}
        comment: "updated: ${{ github.event.comment.body }} from: ${{ github.event.changes.body.from }}"
        
    - name: Delete issue
      id: del_issue
      if: ${{ github.event.action == 'deleted' }}
      uses: atlassian/gajira-comment@master
      with:
        issue: ${{ steps.find.outputs.issue }}
        comment: "deleted: ${{ github.event.comment.body }}"
`
