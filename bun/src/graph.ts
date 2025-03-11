type Degree3Node<T> = {
	value: T;
	neighbors: [
		Degree3Node<T> | null,
		Degree3Node<T> | null,
		Degree3Node<T> | null
	];
};

type ArbitraryNeighborNode<T> = {
	value: T;
	neighbors: (ArbitraryNeighborNode<T> | null)[]
}

// TODO: cache depths
function getDepth<T>(
	node: Degree3Node<T> | null,
	visits: Set<Degree3Node<T>>
): number {
	if (node === null) {
		return 0;
	}

	if (visits.has(node)) {
		return 0;
	}

	visits.add(node);

	return (
		1 +
		Math.max(
			getDepth(node.neighbors[0], visits),
			getDepth(node.neighbors[1], visits),
			getDepth(node.neighbors[2], visits)
		)
	);
}

function getSize<T>(
	predecessor: Degree3Node<T> | null,
	visits: Set<Degree3Node<T>>
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
		getSize(predecessor.neighbors[0], visits) +
		getSize(predecessor.neighbors[1], visits) +
		getSize(predecessor.neighbors[2], visits)
	);
}

function getTreeSize(
	parent: ArbitraryNeighborNode<unknown> | null,
	visits: Set<ArbitraryNeighborNode<unknown>>,
	sizes: Map<ArbitraryNeighborNode<unknown> | null, number>
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

		const size = getTreeSize(neighbor, visits, sizes);

		sizes.set(neighbor, size);

		total += size;
	}
	
	return total;
}

function insert<T>(predecessor: Degree3Node<T>, value: T) {
	// Check the successor spans, and find the one with the least number of nodes.
	let minSizeIndex = 0;
	let minSize = Infinity;
	for (let i = 0; i < predecessor.neighbors.length; i++) {
		let size = getSize(predecessor.neighbors[i] ?? null, new Set());
		if (size < minSize) {
			minSize = size;
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
		insert(child, value);
	}
}

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

// TODO: check for visited nodes
function* breadthFirstSearch<T>(node: ArbitraryNeighborNode<T> | null): Generator<ArbitraryNeighborNode<T>> {
	if (node === null) {
		return;
	}

	const queue = [node];

	while (queue.length > 0) {
		const current = queue.shift();
		if (!current) continue;

		yield current;

		queue.push(...current.neighbors.filter((n) => n !== null));
	}
}

function isTree(predecessor: Degree3Node<unknown> | null, visited: Set<Degree3Node<unknown>>): boolean {
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

function areDisjoint<T>(a: Degree3Node<T> | null, b: Degree3Node<T> | null): boolean {
	if (a === null || b === null) {
		return true;
	}

	const aSet = new Set<Degree3Node<T>>(depthFirstSearch(a, new Set()));
	const bSet = new Set<Degree3Node<T>>(depthFirstSearch(b, new Set()));

	return aSet.size + bSet.size === aSet.union(bSet).size;
}

function firstFreeNode<T>(node: Degree3Node<T> | null): Degree3Node<T> | null {
	if (node === null) {
		return null;
	}

	for (const successor of breadthFirstSearch(node)) {
		if (successor.neighbors.some((n) => n === null)) {
			return successor;
		}
	}

	return null;
}

function centroid<T>(node: ArbitraryNeighborNode<T> | null, visited: Set<ArbitraryNeighborNode<T>>): ArbitraryNeighborNode<T> | null {
	if (node === null) {
		return null;
	}
	
	const sizes = new Map<ArbitraryNeighborNode<T> | null, number>();
	const treeSize = getTreeSize(node, visited, sizes);

	let centroid = node;

	
	
	return centroid;
}

function joinTrees<T>(a: Degree3Node<T> | null, b: Degree3Node<T> | null): Degree3Node<T> | null {
	if (!isTree(a, new Set())) {
		throw new Error("a is not a tree");
	}

	if (!isTree(b, new Set())) {
		throw new Error("b is not a tree");
	}

	if (!areDisjoint(a, b)) {
		throw new Error("a and b are not disjoint");
	}

	if (a === null) {
		return b;
	}

	if (b === null) {
		return a;
	}

	const freeNodeA = firstFreeNode(a);
	if (freeNodeA === null) {
		throw new Error("a has no free nodes");
	}

	const freeNodeB = firstFreeNode(b);
	if (freeNodeB === null) {
		throw new Error("b has no free nodes");
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

	// TODO: find the centroid of the tree `a`.
	
	return a;
}

/**
 * Deletes one of the successors of `predecessor` that satisfies the predicate.
 * 
 * @param predecessor the predecessor node for which to delete one of its
 *   successors.
 * @param predicate the predicate to use to determine which successor to delete.
 */
function deleteNode<T>(predecessor: Degree3Node<T> | null, predicate: (value: T) => boolean) {
	if (predecessor === null) {
		return;
	}

	for (let i = 0; i < predecessor.neighbors.length; i++) {
		const successor = predecessor.neighbors[i] ?? null;
		if (successor === null) {
			continue;
		}

		if (predicate(successor.value)) {
			predecessor.neighbors[i] = null;
			return;
		}
		
		deleteNode(successor, predicate);
	}
}

class Graph<T> {
	private root: Degree3Node<T> | null = null;

	insert(value: T) {
		if (this.root === null) {
			this.root = {
				value: value,
				neighbors: [null, null, null],
			};
			return;
		}

		insert(this.root, value);
	}

	delete(predicate: (value: T) => boolean) {
		if (this.root === null) {
			return;
		}
	}

	static *depthFirstSearch<T>(graph: Graph<T>) {
		for (const node of depthFirstSearch(graph.root, new Set())) {
			yield node.value;
		}
	}

	static *breadthFirstSearch<T>(graph: Graph<T>) {
		for (const node of breadthFirstSearch(graph.root)) {
			yield node.value;
		}
	}
}
