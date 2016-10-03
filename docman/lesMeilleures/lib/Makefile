CXX         =	g++ -std=c++14

MKDIR       =   mkdir -p

RM          =   rm -rf

CXXFLAGS    =	-Wfatal-errors -Werror -Wall -Wextra -pedantic

LDFLAGS     =	

LDLIBS      =	

PRG         =	liblesMeilleures.so

INCS_DIR    =   ./inc/

LIBS_DIR    =   ./lib/

OBJS_DIR    =   ./obj/

SRCS_DIR    =   ./src/

TPLS_DIR    =   ./tpl/

SRCS	    =	$(wildcard $(SRCS_DIR)/*.cpp)      \

OBJS        :=  $(SRCS:$(SRCS_DIR)%.cpp=$(OBJS_DIR)%.o)

DEPS        :=  $(OBJS:$(OBJS_DIR)%.o=$(OBJS_DIR)%.d)

#---RULES---
all: objDir $(PRG)
	@echo -e "\t\033[1;32mBuild done!\033[0m"

objDir:
	@$(MKDIR) $(OBJS_DIR)

-include $(DEPS)

$(PRG): $(OBJS)
	@$(CXX) -shared $(OBJS) -o $@ $(LDFLAGS) $(LDLIBS)
	@echo -e "\033[K\033[1;37mCompilation: \t\t\033[32m✓\033[0m"

$(OBJS_DIR)%.o: $(SRCS_DIR)%.cpp
	@echo -e "\033[K\033[1;37mCompilation: \t\t\033[1;33m...\033[0m"
	@echo -en "\033[1A"
	@$(CXX) -c -fPIC $< -I$(INCS_DIR) -I$(TPLS_DIR) -I$(LIBS_DIR) $(CXXFLAGS) -o $@ || (echo -e "\033[K\033[1;37mCompilation: \t\t\033[1;31m✗\033[0m"; false)
	@$(CXX) -M $< -I$(INCS_DIR) -I$(TPLS_DIR) -I$(LIBS_DIR) -MT $@ -MF $(OBJS_DIR)$*.d

clean:
	@echo -en "\033[1;37mObjects clean:\t\t"
	@$(RM) $(OBJS) $(DEPS) $(OBJS_DIR)
	@echo -e "\033[1;32m✓\033[0m"

distclean: clobber
realclean: clobber
clobber: clean
	@echo -ne "\033[1;37mDirectory clean:\t"
	@$(RM) $(PRG)
	@echo -e "\033[1;32m✓\033[0m"

.PHONY: all clean distclean realclean clobber
