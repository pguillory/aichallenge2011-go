package main

//import "fmt"

type Scrum struct {
    points *PointSet
}

func NewScrum() *Scrum {
    this := new(Scrum)
    this.Reset()
    return this
}

func (this *Scrum) Reset() {
    this.points = new(PointSet)
}

func (this *Scrum) Include(p Point) {
    this.points.Include(p)
}

func (this *Scrum) At(p Point) bool {
    return this.points.Includes(p)
}
