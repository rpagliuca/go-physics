#version 410 core

layout (location = 0) in vec3 position;
layout (location = 1) in float ourGreen;

uniform mat4 world;
uniform mat4 camera;
uniform mat4 project;

out vec4 ourColor;

void main()
{
    gl_Position = project * camera * world * vec4(position, 1.0);
    ourColor = vec4(0.0f, ourGreen, 0.8f, 1.0f);
}
