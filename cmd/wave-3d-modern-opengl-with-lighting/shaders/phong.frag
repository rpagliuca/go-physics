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
	vec3 blue = vec3(15.0f/255.0f, 68.0f/255.0f, 255.0f/255.0f) * 1;

	// ambient
	vec3 ambient = blue * 0.7;

	// diffuse
	vec3 norm = normalize(Normal);
	vec3 dirToLight = normalize(LightPos - FragPos);
	float lightNormalDiff = 0.5 * (dot(norm, dirToLight) + 1);
	vec3 diffuse = vec3(0.0f, 0.3f, 0.7f) * 0.5 * lightNormalDiff;

	vec3 viewPos = vec3(0.0f, 0.0f, 0.0f);

	// specular new (blinn-phong)
    vec3 lightDir = normalize(LightPos - FragPos);
    vec3 viewDir = normalize(viewPos - FragPos);
    vec3 halfwayDir = normalize(lightDir + viewDir);
    float specblinn = pow(max(dot(norm, halfwayDir), 0.0), 16.0);

    // specular old
	vec3 dirToView = normalize(viewPos - FragPos);
	vec3 reflectDir = reflect(-dirToLight, norm);
	float specphong = pow(max(dot(dirToView, reflectDir), 0.0), 16.0);

	//vec3 specular = specphong * vec3(-5.0f, 10.0f, -5.0f) * 1;
	vec3 specular = specphong * vec3(1.0f, 0.3f, -1.0f) * 5;

    // Light attenuation over distance
    float attenuation = 0.00140 * pow(dot(LightPos-FragPos, LightPos-FragPos), 0.5);
    //attenuation = 1;

	vec3 result = (diffuse + specular + ambient) / attenuation;
    //result = vec3(0.0f);
	color = vec4(result, 1.0f);
}
