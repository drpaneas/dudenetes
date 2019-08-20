workflow "Build and deploy on push" {
  on = "push"
  resolves = ["ci"]
}

# Setups: ~/go/src/github.com/drpaneas/dudenetes

action "ci" {
  uses="cedrickring/golang-action@1.3.0"
}
