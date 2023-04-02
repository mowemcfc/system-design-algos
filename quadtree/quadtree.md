# Quadtree

A quadtree is a tree data structure in which each internal node has exactly four children. 

Quadtrees are often used to partition a two-dimensional space by recursively subdividing it into four quadrants or regions. This recursive partitioning occurs until a defined maximum depth, space, or points-per-quadrant constraint is met.

## Motivation

For a given 2D space, imagine a system in which there can be n = (l x w) different defined "points" within the 2D space.

To calculate the difference between one point and all others, you need to perform n^2 operations.

In order to define set of "regions" that can effectively subdivide the space such that you can perform operations related to points inside each region efficiently.

In practice, you only need to subdivide a space once the amount of points in a region crosses some kind of tractability threshold.

By doing so, you can reduce the complexity of the above distance calculation in O(nlogn) time.
