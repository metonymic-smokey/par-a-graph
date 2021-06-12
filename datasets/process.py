import sys

if len(sys.argv) < 3:
    ORIG_FILE = "./wiki-Vote.txt"
    PREFIX = "wiki-vote"
else:
    ORIG_FILE = sys.argv[1]
    PREFIX = sys.argv[2]

print("Using file:", ORIG_FILE)
print("Output prefix:", PREFIX)

nodes = set()

with open(f"{PREFIX}-edges.txt", "w") as edge_file:
    with open(ORIG_FILE) as f:
        for line in f:
            line = line.strip()
            if line.startswith("#") or line == "":
                continue

            source, dest = line.split()
            print(f"{source},{dest}", file=edge_file)
            nodes.add(source)
            nodes.add(dest)

with open(f"{PREFIX}-nodes.txt", "w") as node_file:
    for node in nodes:
        print(f"{node},{node}", file=node_file)
