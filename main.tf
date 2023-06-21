resource "terraform_data" "hello-world" {
  provisioner "local-exec" {
    command = "echo Hello World"
  }
}
