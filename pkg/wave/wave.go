package wave

const LEN = 100
const D = 0.005

type Grid [3][LEN]float64

// Wave equation
//
// d2u/dt2 = c2 d2u/dx2
//
// du/dt (forward): (u(t+dt) - u(t)) / dt
//
// d2u/dt2 (x fixo): ((u(t+dt) - u(t)) / dt - (u(t) - u(t-dt)) / dt) / dt
//				= (u(t+dt) - 2u(t) + u(t-dt)) / dt^2
//
// d2u/dx2 (t fixo): ((u(x+dx) - u(x)) / dx - (u(x) - u(x-dx)) / dx) / dx
//				= (u(x+dx) - 2u(x) + u(x-dx)) / dx^2
//
// Equação combinada:
// (u(t+dt) - 2u(t) + u(t-dt)) / dt^2 = C * (u(x+dx) - 2u(x) + u(x-dx)) / dx^2
//				=> u(t+dt) = D * (u(x+dx) - 2u(x) + u(x-dx)) + 2u(t) - u(t-dt)

func NextStep(grid Grid) Grid {

	for i := 1; i < LEN-1; i++ {
		// Combined equation
		grid[2][i] = D*(grid[1][i+1]-2.0*grid[1][i]+grid[1][i-1]) + 2.0*grid[1][i] - grid[0][i]
	}

	// Sane boundaries
	grid[2][0] = grid[2][1]
	grid[2][LEN-1] = grid[2][LEN-2]

	// Cycle temporal values
	grid[0] = grid[1]
	grid[1] = grid[2]

	return grid
}
