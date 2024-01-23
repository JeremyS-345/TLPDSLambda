data "aws_iam_policy_document" "assume_role" {
  statement {
    effect = "Allow"

    principals {
      type        = "Service"
      identifiers = ["lambda.amazonaws.com"]
    }

    actions = ["sts:AssumeRole"]
  }
}

resource "aws_iam_role" "iam_for_lambda" {
  name               = "iam_for_dynamo_lambda"
  assume_role_policy = data.aws_iam_policy_document.assume_role.json
}

data "archive_file" "lambda" {
  type        = "zip"
  source_file = "../ingestData/bootstrap"
  output_path = "bootstrap.zip"
}

resource "aws_lambda_function" "LPDS_lambda" {
  # If the file is not in the current working directory you will need to include a
  # path.module in the filename.
  filename      = "bootstrap.zip"
  function_name = "LPDS-Lambda"
  role          = aws_iam_role.iam_for_lambda.arn
  handler       = "index.test"

  source_code_hash = data.archive_file.lambda.output_base64sha256

  runtime = "provided.al2"

  environment {
    variables = {
      foo = "bar"
    }
  }
}

resource "aws_iam_role_policy_attachment" "lambda_policy" {
   role = aws_iam_role.iam_for_lambda.name
   policy_arn = aws_iam_policy.lambda_dynamo_policy.arn
}


resource "aws_iam_policy" "lambda_dynamo_policy" {
  name        = "lambda-dynamo-policy"
  policy      = data.aws_iam_policy_document.lambda_dynamo_policy.json
}

data "aws_iam_policy_document" "lambda_dynamo_policy" {
  statement {
    effect    = "Allow"
    actions   = ["dynamodb:*"]
    resources = [aws_dynamodb_table.basic_dynamodb_table.arn]
  }
}