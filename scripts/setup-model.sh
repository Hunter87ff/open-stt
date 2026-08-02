export LLAMA_CACHE=./models
llama-server -hf ggml-org/Qwen3-ASR-0.6B-GGUF:Q8_0 -c 256 -np 1 -b 32 -ub 32 --flash-attn on --n-gpu-layers 999 --no-mmap -t 2