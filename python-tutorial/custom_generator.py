front_nine_scores = [3, 3, 4, 4, 5, 5, 6, 6, 7]
def list_front_nine_scores():
    for i, s in enumerate(front_nine_scores):
        yield "Hole #{hole}: {score}:".format(hole=i, score=s)

for s in list_front_nine_scores():
    print s
