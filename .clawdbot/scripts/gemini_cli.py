#!/usr/bin/env python3
"""Minimal Gemini CLI wrapper used by the orchestrator.

Usage examples:

  # Simple prompt string
  ./gemini_cli.py --model gemini-2.0-pro-exp -p "Explain the architecture"

  # Read prompt from file (recommended for longer context)
  ./gemini_cli.py --prompt-file ../prompts/feat.md

Environment:
  GOOGLE_API_KEY must be set with a valid Gemini API key.
"""

import argparse
import os
import sys

try:
    import google.generativeai as genai
except ModuleNotFoundError:
    print("[gemini_cli] Missing dependency: google-generativeai.\n"
          "Install with: pip install google-generativeai", file=sys.stderr)
    sys.exit(1)


def parse_args() -> argparse.Namespace:
    parser = argparse.ArgumentParser(description="Gemini CLI wrapper")
    parser.add_argument("prompt", nargs="?", help="Prompt string (optional if using --prompt-file)")
    parser.add_argument("prompt_tail", nargs=argparse.REMAINDER,
                        help="Additional prompt fragments (allowing spawn script to pass multi-word prompts)")
    parser.add_argument("--model", default=os.environ.get("GEMINI_MODEL", "gemini-2.0-pro-exp"),
                        help="Gemini model ID (default from env GEMINI_MODEL or gemini-2.0-pro-exp)")
    parser.add_argument("-p", "--prompt-string", dest="prompt_string",
                        help="Explicit prompt string (highest priority)")
    parser.add_argument("-f", "--prompt-file", dest="prompt_file",
                        help="Path to file containing prompt text")
    parser.add_argument("-s", "--system", dest="system",
                        help="Optional system instruction for the model")
    parser.add_argument("--json", action="store_true",
                        help="Print raw JSON response instead of plain text")
    return parser.parse_args()


def load_prompt(args: argparse.Namespace) -> str:
    # Priority: --prompt-string > --prompt-file > positional args > STDIN
    if args.prompt_string:
        return args.prompt_string
    if args.prompt_file:
        try:
            return open(args.prompt_file, "r", encoding="utf-8").read()
        except OSError as exc:
            print(f"[gemini_cli] Unable to read prompt file: {exc}", file=sys.stderr)
            sys.exit(1)
    parts = []
    if args.prompt:
        parts.append(args.prompt)
    if args.prompt_tail:
        parts.extend(args.prompt_tail)
    if parts:
        return " ".join(parts).strip()
    if not sys.stdin.isatty():
        return sys.stdin.read()
    print("[gemini_cli] No prompt provided. Use -p/--prompt-file or pipe input.", file=sys.stderr)
    sys.exit(1)


def main():
    args = parse_args()
    api_key = os.environ.get("GOOGLE_API_KEY")
    if not api_key:
        print("[gemini_cli] GOOGLE_API_KEY not set.", file=sys.stderr)
        sys.exit(1)

    prompt = load_prompt(args)

    genai.configure(api_key=api_key)
    generation_config = {
        "temperature": 0.4,
        "top_p": 0.95,
        "top_k": 32,
        "max_output_tokens": 2000,
    }
    model_kwargs = {"model_name": args.model, "generation_config": generation_config}
    if args.system:
        model_kwargs["system_instruction"] = args.system

    model = genai.GenerativeModel(**model_kwargs)
    response = model.generate_content(prompt)

    if args.json:
        import json
        print(json.dumps(response.to_dict(), indent=2))
    else:
        text = getattr(response, "text", None)
        if text:
            print(text)
        else:
            # Fallback for multi-part responses
            parts = []
            for part in response.candidates[0].content.parts:
                if getattr(part, "text", None):
                    parts.append(part.text)
            print("\n".join(parts))


if __name__ == "__main__":
    main()
