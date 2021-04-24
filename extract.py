

f = 'input.txt'
f2 = open('dirLinks.txt','w')

with open(f) as fp:
    line = fp.readline()
    while line:
        line = line.strip()
        l = line.split()
        l1 = l[0]
        l2 = " ".join(l[1:])
        l = [l1,l2]
        l = ",".join(l) + "\n"
        f2.write(l)
        line = fp.readline()

f2.close()
