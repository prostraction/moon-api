// Renders the lit fraction of the Moon as an SVG path.
//
// Geometry: the terminator is a half-ellipse from pole to pole.
//   - semi-major axis (vertical) = R
//   - semi-minor axis (horizontal) = R * |1 - 2*illum|
// Together with a half-circle along the lit limb, this gives the classic
// crescent/gibbous shape.
//
// We assume Northern-Hemisphere viewing convention:
//   - waxing  → lit side on the right
//   - waning  → lit side on the left

const R = 50;
const VIEWBOX = `${-R} ${-R} ${R * 2} ${R * 2}`;

// Renders the SVG once into the host element, then `update()` recolors and
// reshapes it for the given illumination/waxing values.
export function createMoonShape(host) {
  host.innerHTML = `
    <svg id="moon-svg" viewBox="${VIEWBOX}" overflow="visible"
         xmlns="http://www.w3.org/2000/svg" aria-hidden="true">
      <defs>
        <radialGradient id="litTexture" cx="50%" cy="50%" r="50%">
          <stop offset="0%"  stop-color="#f3eddc"/>
          <stop offset="70%" stop-color="#d8d2bc"/>
          <stop offset="100%" stop-color="#9e9a82"/>
        </radialGradient>
        <!-- Dark side: medium-dark neutral gray. Solid enough to read as
             unlit matter (not translucent), warm/desaturated enough to
             avoid looking "dimly illuminated". -->
        <radialGradient id="darkTexture" cx="50%" cy="50%" r="50%">
          <stop offset="0%"  stop-color="#363639"/>
          <stop offset="65%" stop-color="#2a2a2c"/>
          <stop offset="100%" stop-color="#1c1c1f"/>
        </radialGradient>
        <!-- crater pattern, reused via <use href> twice with different fills.
             Stored as a plain <g> (not <symbol>) so coordinates resolve in
             the parent viewBox without needing explicit width/height. -->
        <g id="craters">
          <circle cx="-15" cy="-20" r="6" />
          <circle cx="15"  cy="-25" r="4" />
          <circle cx="-25" cy="0"   r="5" />
          <circle cx="-5"  cy="5"   r="3" />
          <circle cx="10"  cy="10"  r="6" />
          <circle cx="20"  cy="15"  r="4" />
          <circle cx="-15" cy="25"  r="5" />
          <circle cx="30"  cy="0"   r="3" />
          <circle cx="-30" cy="15"  r="2" />
        </g>
        <!-- Soft halo emitted by the lit area only. Larger filter region so
             the blur isn't clipped against the disc bounds. -->
        <filter id="moonGlow" x="-30%" y="-30%" width="160%" height="160%">
          <feGaussianBlur stdDeviation="1.5" result="halo"/>
          <feComponentTransfer in="halo" result="haloAlpha">
            <feFuncA type="linear" slope="0.35"/>
          </feComponentTransfer>
          <feMerge>
            <feMergeNode in="haloAlpha"/>
            <feMergeNode in="SourceGraphic"/>
          </feMerge>
        </filter>
        <!-- Soft Gaussian blur for crater shapes — turns geometric circles
             into hazy patches that read as natural depressions. -->
        <filter id="craterSoftness" x="-20%" y="-20%" width="140%" height="140%">
          <feGaussianBlur stdDeviation="0.7"/>
        </filter>
        <!-- A clipped circle that traces the lit limb only — used as the
             glow source so the halo follows the crescent/gibbous edge. -->
        <clipPath id="litClip">
          <path id="litPath" d=""/>
        </clipPath>
      </defs>

      <!-- dark side: full disc -->
      <circle cx="0" cy="0" r="${R}" fill="url(#darkTexture)"/>

      <!-- dark-side craters: even darker than the unlit gradient — read as
           deeper pockets of shadow rather than as illuminated features. -->
      <use href="#craters" fill="#000000" opacity="0.55"
           filter="url(#craterSoftness)"/>

      <!-- Halo layer: a slightly oversized disc, blurred, clipped to the lit
           shape. Renders BEFORE the lit texture so the halo bleeds out past
           the limb but the surface itself stays sharp. -->
      <g clip-path="url(#litClip)" filter="url(#moonGlow)">
        <circle cx="0" cy="0" r="${R + 1}" fill="#fff7d6" opacity="0.18"/>
      </g>

      <!-- lit side: clipped to the terminator path -->
      <g clip-path="url(#litClip)">
        <circle cx="0" cy="0" r="${R}" fill="url(#litTexture)"/>
        <use href="#craters" fill="rgba(95,86,68,0.5)"/>
      </g>

      <!-- Lit-side rim — only along the illuminated limb. -->
      <g clip-path="url(#litClip)">
        <circle cx="0" cy="0" r="${R}" fill="none" stroke="rgba(255,255,235,0.35)" stroke-width="0.6"/>
      </g>
    </svg>
  `;
}

// illum: 0..1, isWaxing: bool
export function updateMoonShape(host, illum, isWaxing) {
  const path = host.querySelector('#litPath');
  if (!path) return;

  // Clamp and handle degenerate ends so they don't render as 0-width slivers.
  const i = Math.max(0, Math.min(1, illum));

  if (i <= 0.01) {
    // New moon: nothing lit.
    path.setAttribute('d', '');
    return;
  }
  if (i >= 0.99) {
    // Full moon: lit = full disc.
    path.setAttribute('d', `M ${-R} 0 A ${R} ${R} 0 1 0 ${R} 0 A ${R} ${R} 0 1 0 ${-R} 0 Z`);
    return;
  }

  const b = R * Math.abs(1 - 2 * i); // terminator semi-minor axis
  const gibbous = i > 0.5;

  // SVG arc flags. Pole-to-pole drawing:
  //   start at top pole (0, -R)
  //   draw lit limb (half-circle) down to bottom pole (0, +R)
  //   draw terminator (half-ellipse) back to top pole
  //
  // For NH waxing (lit on the right):
  //   limb arc sweeps clockwise → sweep=1
  //   terminator arc sweeps counter-clockwise for crescent (sweep=0),
  //   clockwise for gibbous (sweep=1)
  //
  // For waning, mirror by negating x of both arcs (limb sweeps the other way).

  const limbSweep = isWaxing ? 1 : 0;
  const termSweep = (isWaxing ? 0 : 1) ^ (gibbous ? 1 : 0);

  const d = [
    `M 0 ${-R}`,
    `A ${R} ${R} 0 0 ${limbSweep} 0 ${R}`,
    `A ${b} ${R} 0 0 ${termSweep} 0 ${-R}`,
    'Z',
  ].join(' ');

  path.setAttribute('d', d);
}
