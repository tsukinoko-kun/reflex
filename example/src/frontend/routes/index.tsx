import { useState } from "react"

export default function Page() {
	const [count, setCount] = useState(0)
	return (
		<div className="mx-auto block max-w-4xl">
			<h1 className="text-4xl text-red-700">Hello, Reflex!</h1>
			<p className="text-pretty">Welcome to the Reflex example app.</p>
			<p className="text-pretty">
				Lorem ipsum dolor sit amet, consetetur sadipscing elitr, sed diam nonumy eirmod tempor invidunt ut labore et
				dolore magna aliquyam erat, sed diam voluptua. At vero eos et accusam et justo duo dolores et ea rebum. Stet
				clita kasd gubergren, no sea takimata sanctus est Lorem ipsum dolor sit amet. Lorem ipsum dolor sit amet,
				consetetur sadipscing elitr, sed diam nonumy eirmod tempor invidunt ut labore et dolore magna aliquyam erat, sed
				diam voluptua. At vero eos et accusam et justo duo dolores et ea rebum. Stet clita kasd gubergren, no sea
				takimata sanctus est Lorem ipsum dolor sit amet.
			</p>
			<div className="flex flex-col items-center">
				<button
					type="button"
					onClick={() => setCount(count + 1)}
					className="cursor-pointer font-bold text-2xl text-blue-500"
				>
					+
				</button>
				<p>Count: {count}</p>
				<button
					type="button"
					onClick={() => setCount(count - 1)}
					className="cursor-pointer font-bold text-2xl text-blue-500"
				>
					-
				</button>
			</div>
		</div>
	)
}
