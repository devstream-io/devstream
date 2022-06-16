# Contributor Ladder

## Contributor Ladder

Hello! We are excited that you want to learn more about our project contributor ladder! This contributor ladder outlines the different contributor roles within the project, along with the responsibilities and privileges that come with them. Community members generally start at the first levels of the "ladder" and advance up it as their involvement in the project grows.  Our project members are happy to help you advance along the contributor ladder.

Each of the contributor roles below is organized into lists of three types of things. "Responsibilities" are things that contributor is expected to do. "Requirements" are qualifications a person needs to meet to be in that role, and "Privileges" are things contributors on that level are entitled to.

### Contributor

_Description: A Contributor contributes directly to the project and adds value to it. Contributions need not be code. Non-code contributions are equally, if not more, valued. People at the Contributor level may be new contributors, or they may only contribute occasionally._

- Responsibilities:
    - Follow the CNCF CoC
    - Follow the project contributing guide

- Requirements (one or several of the below items):
    - Report and sometimes resolve issues
    - Occasionally submit PRs
    - Contribute to the documentation
    - Show up at meetings, takes notes
    - Answer questions from other community members
    - Submit feedback on issues and PRs
    - Test releases and patches and submit reviews
    - Run or helps run events
    - Promote the project in public
    - Help run the project infrastructure

- Privileges:
    - Invitations to contributor events
    - Eligible to become an Organization Member

### Organization Member

_Description: An Organization Member is an established contributor who regularly participates in the project. Organization Members have privileges in both project repositories and elections, and as such are expected to act in the interests of the whole project._

An Organization Member _must meet the responsibilities and has the requirements of a Contributor_, plus:

- Responsibilities:
    - Continues to contribute regularly, as demonstrated by having at least 50 GitHub contributions per year, as demonstrated by [GitHub Insights](https://github.com/devstream-io/devstream/pulse).

- Requirements:
    - Must have successful contributions to the project, including at least one of the following:
        - 5 accepted PRs,
        - Reviewed 5 PRs,
        - Resolved and closed 5 Issues,
        - Become responsible for a key project management area,
        - Or some equivalent combination or contribution
    - Must have been contributing for at least 1 months
    - Must be actively contributing to at least one project area
    - Must have two sponsors who are also Organization Members, at least one of whom does not work for the same employer

- Privileges:
    - May be assigned Issues and Reviews
    - May give commands to CI/CD automation
    - Entitled to vote in elections [TODO]
    - Can be added to @devstream teams
    - Can recommend other contributors to become Org Members

The process for a Contributor to become an Organization Member is as follows:

1. The member contributor is nominated by a sponsor, by opening an issue in the appropriate repository.
2. The second sponsor from a different employer seconds the nomination in the issue.
3. The nomination is reviewed publicly in a community meeting. All members vote for the nomination. A majority is required in order to approve the nomination.
4. Publicly announce the result in the issue, and in the contributor slack channel; and invite the member into the member channel.

### Reviewer

_Description: A Reviewer has responsibility for specific code, documentation, test, or other project areas. They are collectively responsible, with other Reviewers, for reviewing all changes to those areas and indicating whether those changes are ready to merge. They have a track record of contribution and review in the project._

Reviewers are responsible for a "specific area." This can be a specific code directory, driver, chapter of the docs, test job, event, or other clearly-defined project component that is smaller than an entire repository or subproject. Most often it is one or a set of directories in one or more Git repositories. The "specific area" below refers to this area of responsibility.

Reviewers have all the rights and responsibilities of an Organization Member, plus:

- Responsibilities:
    - Following the reviewing guide
    - Reviewing most Pull Requests against their specific areas of responsibility
    - Reviewing at least 20 PRs per year
    - Helping other contributors become reviewers

- Requirements:
    - Experience as a Contributor for at least 1 months
    - Is an Organization Member
    - Has reviewed, or helped review, at least 5 Pull Requests
    - Has analyzed and resolved test failures in their specific area
    - Has demonstrated an in-depth knowledge of the specific area
    - Commits to being responsible for that specific area
    - Is supportive of new and occasional contributors and helps get useful PRs in shape to commit

- Additional privileges:
    - Has GitHub or CI/CD rights to approve pull requests in specific directories
    - Can recommend and review other contributors to become Reviewers
    - Is listed as Approver in the OWNERS file for certain directories

The process of becoming a Reviewer is:

1. The contributor is nominated by opening a PR against the appropriate repository, which adds their GitHub username to the OWNERS file for one or more directories.
2. At least two members of the team that owns that repository or main directory, who are already Approvers, approve the PR.

### Maintainer

_Description: Maintainers are very established contributors who are responsible for the entire project. As such, they have the ability to approve PRs against any area of the project, and are expected to participate in making decisions about the strategy and priorities of the project._

A Maintainer must meet the responsibilities and requirements of a Reviewer, plus:

- Responsibilities:
    - Reviewing at least 40 PRs per year, especially PRs that involve multiple parts of the project
    - Mentoring new Reviewers
    - Writing refactoring PRs
    - Participating in CNCF maintainer activities
    - Determining strategy and policy for the project
    - Participating in, and leading, community meetings

- Requirements
    - Experience as a Reviewer for at least 3 months
    - Demonstrates a broad knowledge of the project across multiple areas
    - Is able to exercise judgement for the good of the project, independent of their employer, friends, or team
    - Mentors other contributors
    - Can commit to spending at least 40 hours per month working on the project

- Additional privileges:
    - Approve PRs to any area of the project
    - Represent the project in public as a Maintainer
    - Communicate with the CNCF on behalf of the project
    - Have a vote in Maintainer decision-making meetings
    
Process of becoming a maintainer:

1. Any current Maintainer may nominate a current Reviewer to become a new Maintainer, by opening a PR against the root of the DevStream repo, adding the nominee as an Approver in the OWNERS file.
2. The nominee will add a comment to the PR testifying that they agree to all requirements of becoming a Maintainer.
3. A majority of the current Maintainers must then approve the PR.

## Inactivity

It is important for contributors to be and stay active to set an example and show commitment to the project. Inactivity is harmful to the project as it may lead to unexpected delays, contributor attrition, and a lost of trust in the project.

- Inactivity is measured by:
    - Periods of no contributions for longer than 2 months
    - Periods of no communication for longer than 1 month
- Consequences of being inactive include:
    - Involuntary removal or demotion
    - Being asked to move to Emeritus status

## Involuntary Removal or Demotion

Involuntary removal/demotion of a contributor happens when responsibilities and requirements aren't being met. This may include repeated patterns of inactivity, extended period of inactivity, a period of failing to meet the requirements of your role, and/or a violation of the Code of Conduct. This process is important because it protects the community and its deliverables while also opens up opportunities for new contributors to step in.

Involuntary removal or demotion is handled through a vote by a majority of the current Maintainers.

## Stepping Down/Emeritus Process

If and when contributors' commitment levels change, contributors can consider stepping down (moving down the contributor ladder) vs moving to emeritus status (completely stepping away from the project).

Contact the Maintainers about changing to Emeritus status, or reducing your contributor level.

## Contact

For inquiries, please reach out to the @pmc chair @IronCore864.
