type Degree3Node<T> = {
	value: T;
	neighbors: [
		Degree3Node<T> | null,
		Degree3Node<T> | null,
		Degree3Node<T> | null
	];
};

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

// TODO: cache sizes
function getSize<T>(
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
		getSize(node.neighbors[0], visits) +
		getSize(node.neighbors[1], visits) +
		getSize(node.neighbors[2], visits)
	);
}

function insert<T>(node: Degree3Node<T>, value: T) {
	// Check the successor spans, and find the one with the least number of nodes.
	let minSizeIndex = 0;
	let minSize = Infinity;
	for (let i = 0; i < node.neighbors.length; i++) {
		let size = getSize(node.neighbors[i] ?? null, new Set());
		if (size < minSize) {
			minSize = size;
			minSizeIndex = i;
		}
	}

	const child = node.neighbors[minSizeIndex] ?? null;
	if (child === null) {
		node.neighbors[minSizeIndex] = {
			value,
			neighbors: [node, null, null],
		};
	} else {
		insert(child, value);
	}
}

function* depthFirstSearch<T>(
	node: Degree3Node<T> | null,
	visited: Set<Degree3Node<T>> = new Set()
): Generator<T> {
	if (node === null) {
		return;
	}

	if (visited.has(node)) {
		return;
	}

	visited.add(node);

	yield node.value;

	for (const neighbor of node.neighbors) {
		yield* depthFirstSearch(neighbor);
	}
}

function* breadthFirstSearch<T>(node: Degree3Node<T> | null): Generator<T> {
	if (node === null) {
		return;
	}

	const queue = [node];

	while (queue.length > 0) {
		const current = queue.shift();
		if (!current) continue;

		yield current.value;

		queue.push(...current.neighbors.filter((n) => n !== null));
	}
}

class Graph<T> {
	private root: Degree3Node<{ id: string; value: T }> | null = null;

	insert(node: { id: string; value: T }) {
		if (this.root === null) {
			this.root = {
				value: node,
				neighbors: [null, null, null],
			};
			return;
		}

		insert(this.root, node);
	}

	static *depthFirstSearch<T>(graph: Graph<T>) {
		for (const value of depthFirstSearch(graph.root)) {
			yield { ...value };
		}
	}

	static *breadthFirstSearch<T>(graph: Graph<T>) {
		for (const value of breadthFirstSearch(graph.root)) {
			yield { ...value };
		}
	}
}
