import random
import string

def generate_random_string(length):
    letters = string.ascii_letters + string.digits
    result_str = ''.join(random.choice(letters) for i in range(length))
    return result_str

# Calculate the number of characters needed for 10KB
char_count = 5120  # 10KB = 10240 bytes

# Generate the random string
random_string = generate_random_string(char_count)

# Print the generated string
print(random_string)