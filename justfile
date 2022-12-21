g:
  bazel run //:gazelle

update-go-repos:
  bazel run //:gazelle update-repos -- -from_file=go.mod

dev: 
  ibazel run weather --run_output_interactive=false