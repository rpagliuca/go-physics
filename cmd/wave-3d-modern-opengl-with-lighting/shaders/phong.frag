#version 410 core

struct Material {
	vec3 ambient;
	vec3 diffuse;
	vec3 specular;
	float shininess;
};

struct Light {
	vec3 ambient;
	vec3 diffuse;
	vec3 specular;
};

in vec3 Normal;
in vec3 FragPos;
in vec3 LightPos;
out vec4 color;

uniform Material material;
uniform Light light;

void main()
{
    vec3 orange = vec3(253.0f/255.0f, 106.0f/255.0f, 2.0f/255.0f);
    vec3 white = vec3(1.0f, 1.0f, 1.0f);
    vec3 darkblue = vec3(0.0f, 0.1f, 0.5f);
    vec3 blue = vec3(0.0f, 0.3f, 0.7f);

	// ambient
	vec3 ambient = darkblue;

	// diffuse
	vec3 norm = normalize(Normal);
	vec3 dirToLight = normalize(FragPos - LightPos);
	float lightNormalDiff = max(dot(norm, dirToLight), 0.0);
	//vec3 diffuse = vec3(0.5f, 1.0f, 0.5f) * lightNormalDiff;
	vec3 diffuse = blue * lightNormalDiff;

	// specular
	vec3 viewPos = vec3(0.0f, 0.0f, 0.0f);
	vec3 dirToView = normalize(viewPos - FragPos);
	vec3 reflectDir = reflect(dirToLight, norm);
	float spec = pow(max(dot(dirToView, reflectDir), 0.0), 64);
	vec3 specular = spec * white;

    // Light attenuation over distance
    float attenuation = 0.0015 * pow(dot(LightPos-FragPos, LightPos-FragPos), 0.5);
    //attenuation = 1;

	vec3 result = (diffuse + specular + ambient) / attenuation;
    //result = vec3(0.0f);
	color = vec4(result, 1.0f);
}
