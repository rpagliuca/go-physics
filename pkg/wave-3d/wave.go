package wave3d

const LEN = 20
const D = 0.06

type Grid [3][LEN][LEN]float64

// Wave equation
//
// d2u/dt2 = c2 (d2u/dx2 + d2u/dy2)
//
// du/dt (forward): (u(t+dt) - u(t)) / dt
//
// d2u/dt2 (x fixo, y fixo): ((u(t+dt) - u(t)) / dt - (u(t) - u(t-dt)) / dt) / dt
//				= (u(t+dt) - 2u(t) + u(t-dt)) / dt^2
//
// d2u/dx2 (t fixo, x fixo): ((u(x+dx) - u(x)) / dx - (u(x) - u(x-dx)) / dx) / dx
//				= (u(x+dx) - 2u(x) + u(x-dx)) / dx^2
//
// d2u/dy2 (t fixo, y fixo): ((u(y+dy) - u(y)) / dy - (u(y) - u(y-dy)) / dy) / dy
//				= (u(y+dy) - 2u(y) + u(y-dy)) / dy^2
//
// Equação combinada:
// (u(t+dt) - 2u(t) + u(t-dt)) / dt^2 = C * (u(x+dx) - 2u(x) + u(x-dx) + u(y+dy) - 2u(y) + u(y-dy)) / dx^2
//				=> u(t+dt) = D * (u(x+dx) - 2u(x) + u(x-dx) + u(y+dy) - 2u(y) + u(y-dy)) + 2u(t) - u(t-dt)

func NextStep(grid Grid) Grid {

	for i := 1; i < LEN-1; i++ {
		for j := 1; j < LEN-1; j++ {
			// Combined equation
			grid[2][i][j] =
				D*
					(grid[1][i+1][j]-2.0*grid[1][i][j]+grid[1][i-1][j]) +
					D*
						(grid[1][i][j+1]-2.0*grid[1][i][j]+grid[1][i][j-1]) +
					2.0*grid[1][i][j] - grid[0][i][j]
		}
	}

	for i := 0; i < LEN; i++ {
		// Sane boundaries
		grid[2][i][0] = grid[2][i][1]
		grid[2][i][LEN-1] = grid[2][i][LEN-2]
	}

	for j := 0; j < LEN; j++ {
		// Sane boundaries
		grid[2][0][j] = grid[2][1][j]
		grid[2][LEN-1][j] = grid[2][LEN-2][j]
	}

	// Cycle temporal values
	grid[0] = grid[1]
	grid[1] = grid[2]

	return grid
}
