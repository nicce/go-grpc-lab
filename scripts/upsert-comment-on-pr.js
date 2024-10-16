/**
 * Applies comment to the PR. If a comment exist, it will be updated.
 * See comment at bottom of this file how this can be used.
 *
 * @param {octokit/rest.js} github      Authenticated GitHub client. See more here:
 *                                      https://octokit.github.io/rest.js/v20
 * @param {object} context              GitHub context of the workflow run. See more here:
 *                                      https://github.com/actions/toolkit/blob/main/packages/github/src/context.ts
 * @param {string} author               GitHub username of the author of the comment.
 * @param {string} header               Header of the comment. Used to uniquely identify the comment to be updated.
 * @param {string} content              Content to add to the PR.
 */
module.exports = async (github, context, author, header, content) => {
    const { data: comments } = await github.rest.issues.listComments({
        owner: context.repo.owner,
        repo: context.repo.repo,
        issue_number: context.issue.number
    });

    const existingComment = comments.find((c) => c.user.login === author && c.body.startsWith(header));

    const actualContent = `${header}\n\n${content}`

    if (existingComment) {
        github.rest.issues.updateComment({
            owner: context.repo.owner,
            repo: context.repo.repo,
            comment_id: existingComment.id,
            body: actualContent
        });
    } else {
        github.rest.issues.createComment({
            owner: context.repo.owner,
            repo: context.repo.repo,
            issue_number: context.issue.number,
            body: actualContent
        });
    }
}

/*
* Example usage:
jobs:
  testing:
    environment: ${{ inputs.environment }}
    runs-on: ubuntu-latest
    steps:
      - name: Checkout code
        uses: actions/checkout@b4ffde65f46336ab88eb53be808477a3936bae11 # v4.1.1

      - name: Add comment on PR
        uses: actions/github-script@60a0d83039c74a4aee543508d2ffcb1c3799cdea #v7.0.1
        with:
          github-token: ${{ secrets.ACCESS_TOKEN }}
          script: |
            const now = new Date();
            const time = now.getHours() + ":" + now.getMinutes() + ":" + now.getSeconds();
            const author = "l-gba-dev-a-itsehbg"
            const content = `
            ## This is the PR comment 🚀 ${time}
            \`\`\`terraform
            Bacon ipsum dolor amet pork loin capicola cow filet mignon chicken pastrami
            \`\`\`
            `
            const header = "# Lorem ipsum"
            const script = require('./scripts/upsert-comment-on-pr.js')
            script(github, context, author, header, content)
*/
