/**
 * A node that will only allow at most 3 neighbours.
 */
type Degree3Node<T> = {
	value: T;
	neighbors: [
		Degree3Node<T> | null,
		Degree3Node<T> | null,
		Degree3Node<T> | null
	];
};

/**
 * Measures the depth of a subgraph.
 *
 * Could be useful for the purposes of finding a neighbouring subgraph that
 * could be a candidate for inserting a new node into.
 * @param predecessor The node for which to measure its depth.
 * @param visits A set of nodes that have already been visited.
 * @returns The depth of the node.
 */
function depth(
	predecessor: Degree3Node<unknown> | null,
	visits: Set<Degree3Node<unknown>>
): number {
	if (predecessor === null) {
		return 0;
	}

	if (visits.has(predecessor)) {
		return 0;
	}

	visits.add(predecessor);

	return (
		1 +
		Math.max(
			depth(predecessor.neighbors[0], visits),
			depth(predecessor.neighbors[1], visits),
			depth(predecessor.neighbors[2], visits)
		)
	);
}

/**
 * Measures the size of a subgraph.
 *
 * Could be useful for the purposes of finding a neighbouring subgraph that
 * could be a candidate for inserting a new node into.
 * @param predecessor The predecessor node for which to measure its size.
 * @param visits A set of nodes that have already been visited.
 * @returns The size of the node.
 */
function size(
	predecessor: Degree3Node<unknown> | null,
	visits: Set<Degree3Node<unknown>>
): number {
	if (predecessor === null) {
		return 0;
	}

	if (visits.has(predecessor)) {
		return 0;
	}

	visits.add(predecessor);

	return (
		1 +
		size(predecessor.neighbors[0], visits) +
		size(predecessor.neighbors[1], visits) +
		size(predecessor.neighbors[2], visits)
	);
}

/**
 * Measures the size of a tree.
 * @param parent The parent node for which to measure its size.
 * @param visits A set of nodes that have already been visited.
 * @param sizes A map of nodes to their sizes.
 * @returns The size of the node.
 */
function treeSize(
	parent: Degree3Node<unknown> | null,
	visits: Set<Degree3Node<unknown>>,
	sizes: Map<Degree3Node<unknown> | null, number>
): number {
	if (parent === null) {
		return 0;
	}

	visits.add(parent);

	let total = 1;
	for (const neighbor of parent.neighbors) {
		if (neighbor === null || neighbor === parent) {
			continue;
		}

		const size = treeSize(neighbor, visits, sizes);

		sizes.set(neighbor, size);

		total += size;
	}

	return total;
}

/**
 * Inserts a value into a degree-3-bounded tree, starting at the specified
 * predecessor node.
 *
 * This algorithm will most certainly not work with an empty tree.
 * @param predecessor The predecessor node for which to insert the value.
 * @param value The value to insert.
 */
function insertIntoIntoAvailableNodeOfDegree3BoundedTree<T>(
	predecessor: Degree3Node<T>,
	value: T
) {
	// Check the successor spans, and find the one with the least number of nodes.
	let minSizeIndex = 0;
	let minSize = Infinity;
	for (let i = 0; i < predecessor.neighbors.length; i++) {
		let sz = size(predecessor.neighbors[i] ?? null, new Set());
		if (sz < minSize) {
			minSize = sz;
			minSizeIndex = i;
		}
	}

	const child = predecessor.neighbors[minSizeIndex] ?? null;
	if (child === null) {
		predecessor.neighbors[minSizeIndex] = {
			value,
			neighbors: [predecessor, null, null],
		};
	} else {
		insertIntoIntoAvailableNodeOfDegree3BoundedTree(child, value);
	}
}

/**
 * Performs a depth-first search, starting at the specified node.
 * @param node The node to start the search from.
 * @param visited A set of nodes that have already been visited.
 * @returns A generator of the nodes in the graph.
 */
function* depthFirstSearch<T>(
	node: Degree3Node<T> | null,
	visited: Set<Degree3Node<T>>
): Generator<Degree3Node<T>> {
	if (node === null) {
		return;
	}

	if (visited.has(node)) {
		return;
	}

	visited.add(node);

	yield node;

	for (const neighbor of node.neighbors) {
		yield* depthFirstSearch(neighbor, visited);
	}
}

/**
 * Performs a breadth-first search, starting at the specified node.
 * @param node The node to start the search from.
 * @param visited A set of nodes that have already been visited.
 * @returns A generator of the nodes in the graph.
 */
function* breadthFirstSearch<T>(
	node: Degree3Node<T> | null,
	visited: Set<Degree3Node<T>>
): Generator<Degree3Node<T>> {
	if (node === null) {
		return;
	}

	const queue = [node];

	while (queue.length > 0) {
		const current = queue.shift();
		if (!current) continue;

		if (visited.has(current)) {
			continue;
		}

		visited.add(current);

		yield current;

		queue.push(...current.neighbors.filter((n) => n !== null));
	}
}

/**
 * Checks if a graph is a tree.
 * @param predecessor The predecessor node for which to check if it is a tree.
 * @param visited A set of nodes that have already been visited.
 * @returns True if the node is a tree, false otherwise.
 */
function isTree(
	predecessor: Degree3Node<unknown> | null,
	visited: Set<Degree3Node<unknown>>
): boolean {
	// Vacuously true
	if (predecessor === null) {
		return true;
	}

	if (visited.has(predecessor)) {
		return false;
	}

	visited.add(predecessor);

	for (const neighbor of predecessor.neighbors) {
		if (neighbor === null) {
			continue;
		}

		if (visited.has(neighbor)) {
			return false;
		}

		if (!isTree(neighbor, visited)) {
			return false;
		}
	}

	return true;
}

/**
 * Checks if the graphs are disjoint.
 * @param a The first graph.
 * @param b The second graph.
 * @returns True if the graphs are disjoint, false otherwise.
 */
function areDisjoint<T>(
	a: Degree3Node<T> | null,
	b: Degree3Node<T> | null
): boolean {
	if (a === null || b === null) {
		return true;
	}

	const aSet = new Set<Degree3Node<T>>(depthFirstSearch(a, new Set()));
	const bSet = new Set<Degree3Node<T>>(depthFirstSearch(b, new Set()));

	return aSet.size + bSet.size === aSet.union(bSet).size;
}

/**
 * Finds the first free node in a tree.
 * @param node The node to search for a free node.
 * @returns The first free node in the tree, or null if there are no free nodes.
 */
function firstFreeNode<T>(node: Degree3Node<T> | null): Degree3Node<T> | null {
	if (node === null) {
		return null;
	}

	for (const successor of breadthFirstSearch(node, new Set())) {
		if (successor.neighbors.some((n) => n === null)) {
			return successor;
		}
	}

	return null;
}

/**
 * Finds the centroid of a tree.
 * @param node The node to find the centroid of.
 * @returns The centroid of the tree, or null if the tree is empty.
 */
function centroid<T>(node: Degree3Node<T> | null): Degree3Node<T> | null {
	if (node === null) {
		return null;
	}

	const sizes = new Map<Degree3Node<T> | null, number>();
	const tSize = treeSize(node, new Set(), sizes);

	for (const [node, size] of sizes.entries()) {
		if (size <= tSize / 2) {
			return node;
		}
	}

	throw new Error("No centroid found. Something went terribly wrong.");
}

/**
 * Joins two trees.
 * @param a The first tree. (If a is not a tree, an error will be thrown.)
 * @param b The second tree. (If b is not a tree, an error will be thrown.)
 * @returns The joined tree.
 */
function joinPairsOfTrees<T>(
	a: Degree3Node<T> | null,
	b: Degree3Node<T> | null
): Degree3Node<T> | null {
	if (!isTree(a, new Set())) {
		throw new Error("a is not a tree");
	}

	if (!isTree(b, new Set())) {
		throw new Error("b is not a tree");
	}

	if (!areDisjoint(a, b)) {
		throw new Error(
			"a and b are not disjoint; likely one is a subtree of the other."
		);
	}

	if (a === null) {
		return b;
	}

	if (b === null) {
		return a;
	}

	const freeNodeA = firstFreeNode(a);
	if (freeNodeA === null) {
		throw new Error("a has no free nodes. Something went terribly wrong.");
	}

	const freeNodeB = firstFreeNode(b);
	if (freeNodeB === null) {
		throw new Error("b has no free nodes. Something went terribly wrong.");
	}

	for (let i = 0; i < freeNodeA.neighbors.length; i++) {
		if (freeNodeA.neighbors[i] === null) {
			freeNodeA.neighbors[i] = freeNodeB;
			break;
		}
	}

	for (let i = 0; i < freeNodeB.neighbors.length; i++) {
		if (freeNodeB.neighbors[i] === null) {
			freeNodeB.neighbors[i] = freeNodeA;
			break;
		}
	}

	a = centroid(a);

	return a;
}

/**
 * Joins a list of trees into a single tree.
 * @param trees The list of trees to join.
 * @returns The joined tree.
 */
function joinTrees<T>(trees: (Degree3Node<T> | null)[]): Degree3Node<T> | null {
	return trees.reduce((a, b) => joinPairsOfTrees(a, b));
}

/**
 * Clears a specific neighbor from a node.
 * @param node The node to clear the neighbor from.
 * @param neighbor The neighbor to clear.
 */
function clearSpecificNeighbor(
	node: Degree3Node<unknown>,
	neighbor: Degree3Node<unknown>
) {
	for (let i = 0; i < node.neighbors.length; i++) {
		if (node.neighbors[i] === neighbor) {
			node.neighbors[i] = null;
		}
	}
}

/**
 * Given the predicate, deletes all successors of `predecessor` that satisfy
 * the predicate.
 * @param predecessor the predecessor node for which to delete one of its
 *   successors.
 * @param predicate the predicate to use to determine which successor to delete.
 */
function deleteNode<T>(
	predecessor: Degree3Node<T> | null,
	predicate: (value: T) => boolean
) {
	if (predecessor === null) {
		return;
	}

	for (const neighbor of predecessor.neighbors) {
		if (neighbor === null) {
			continue;
		}

		if (predicate(neighbor.value)) {
			clearSpecificNeighbor(neighbor, predecessor);
			for (const neighborsNeighbor of neighbor.neighbors) {
				if (neighborsNeighbor === null) {
					continue;
				}

				clearSpecificNeighbor(neighborsNeighbor, neighbor);
			}
			return;
		}
	}
}

/**
 * A class encapsulating methods for manipulating a degree-3-bounded tree.
 */
class Degree3BoundedTree<T> {
	private root: Degree3Node<T> | null = null;

	insertIntoAvailableNode(value: T) {
		if (this.root === null) {
			this.root = {
				value: value,
				neighbors: [null, null, null],
			};
			return;
		}

		insertIntoIntoAvailableNodeOfDegree3BoundedTree(this.root, value);
	}

	/**
	 *
	 * @param predicate Given a predicate, any node who's value satisfies the
	 *   predicate will be deleted.
	 */
	deleteAndNormalizeToCentroid(predicate: (value: T) => boolean) {
		if (this.root === null) {
			return;
		}

		deleteNode(this.root, predicate);

		this.root = centroid(this.root);
	}

	static *depthFirstSearch<T>(graph: Degree3BoundedTree<T>) {
		for (const node of depthFirstSearch(graph.root, new Set())) {
			yield node.value;
		}
	}

	static *breadthFirstSearch<T>(graph: Degree3BoundedTree<T>) {
		for (const node of breadthFirstSearch(graph.root, new Set())) {
			yield node.value;
		}
	}
}
