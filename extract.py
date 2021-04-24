# f = "./input-pageNum.txt"
# f2 = open("pageNum.txt", "w")

# with open(f) as fp:
#     for line in fp:
#         if not line.startswith("#"):
#             line = line.strip()

#             if len(line) == 0:
#                 continue

#             l = line.split()
#             l1 = l[0]
#             l2 = " ".join(l[1:])
#             l = [l1, l2]
#             l = ",".join(l) + "\n"
#             f2.write(l)

# f2.close()

import sys
inp_file = sys.argv[1]
out_file = sys.argv[2]

out = open(out_file, "w")

with open(inp_file) as inp:
    for line in inp:
        line = line.strip()
        if (not line.startswith("#")) and (line != ''):
            l = line.split()
            l1 = l[0]
            l2 = " ".join(l[1:])
            l = [l1, l2]
            l = ",".join(l) + "\n"
            out.write(l)
out.close()