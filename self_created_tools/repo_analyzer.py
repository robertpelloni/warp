import os
import json
import sys

def analyze_repo(path):
    print(f"Analyzing repository at {path}...")
    stats = {
        'total_files': 0,
        'rust_files': 0,
        'go_files': 0,
        'cs_files': 0,
        'java_files': 0,
        'ts_files': 0,
        'js_files': 0,
        'py_files': 0,
    }

    for root, dirs, files in os.walk(path):
        if '.git' in dirs:
            dirs.remove('.git')
        for file in files:
            stats['total_files'] += 1
            if file.endswith('.rs'): stats['rust_files'] += 1
            elif file.endswith('.go'): stats['go_files'] += 1
            elif file.endswith('.cs'): stats['cs_files'] += 1
            elif file.endswith('.java'): stats['java_files'] += 1
            elif file.endswith('.ts') or file.endswith('.tsx'): stats['ts_files'] += 1
            elif file.endswith('.js') or file.endswith('.jsx'): stats['js_files'] += 1
            elif file.endswith('.py'): stats['py_files'] += 1

    print(json.dumps(stats, indent=2))

if __name__ == "__main__":
    analyze_repo(sys.argv[1])
